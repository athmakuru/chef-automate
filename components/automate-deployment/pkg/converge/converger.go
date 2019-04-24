package converge

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	api "github.com/chef/automate/api/interservice/deployment"
	"github.com/chef/automate/components/automate-deployment/pkg/depot"
	"github.com/chef/automate/components/automate-deployment/pkg/habpkg"
	"github.com/chef/automate/components/automate-deployment/pkg/target"
)

const (
	deploymentServiceName    = "deployment-service"
	postgresqlServiceName    = "automate-postgresql"
	pgGatewayServiceName     = "automate-pg-gateway"
	esGatewayServiceName     = "automate-es-gateway"
	elasticsearchServiceName = "automate-elasticsearch"
)

var (
	waitingForReconfigureTimeout = 60 * time.Second
	waitingForRestartTimeout     = 600 * time.Second
)

// A Converger mutates the underlying system in response to system
// events and user commands.
// TODO(ssd) 2019-04-22: new name??
type Converger interface {
	Converge(*Task, DesiredState, EventSink) error
	StopServices(*Task, target.Target, EventSink) error

	// TODO(ssd) 2019-04-21: This is a bit odd. It is necessary
	// for the moment because the actual deployment-service
	// "Reconfigure" happens outside of the Converger.
	ReconfigureDone(*Task) error

	// Stop() stops the Converger, not the deployment-service or
	// any other Automate services.
	Stop()
}

// Converger is a server that handles converge requests. It will take each request,
// compile it into a plan, and execute that plan.
type converger struct {
	stop      chan struct{}
	inbox     chan message
	debug     chan debugReq
	compiler  Compiler
	selfHuper SelfHuper

	currentState state
}

// ConvergerOpt represents a configuration for converger
type ConvergerOpt func(*converger)

// Task is a future like object that will receive a message once
// the requested task has been completed
type Task struct {
	ID uuid.UUID
	// C will receive a message when the task is done
	C chan error
}

// NewTask returns a Task which can be passed to the Converge()
// function. The Task is used to identify and signal completion of
// the requested converge.
func NewTask() (*Task, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "could not create UUID")
	}
	doneChan := make(chan error, 1)
	return &Task{
		ID: id,
		C:  doneChan,
	}, nil
}

// WithCompiler sets the compiler to use when fulfilling a converge request.
// By default, we will use the PhaseOrderedCompiler. Specifying a compiler
// is useful in unit testing
func WithCompiler(compiler Compiler) ConvergerOpt {
	return func(converger *converger) {
		converger.compiler = compiler
	}
}

// WithSelfHuper sets the compiler to use when fulfilling a converge request.
// By default, we will use the PhaseOrderedCompiler. Specifying a compiler
// is useful in unit testing
func WithSelfHuper(s SelfHuper) ConvergerOpt {
	return func(converger *converger) {
		converger.selfHuper = s
	}
}

// WithDebugChannel enables the debug channel used in testing to poll
// the state of the converger. By default this is a nil channel.
//
// NOTE(ssd) 2019-04-22: Ideally I think these messages would also be
// sent via the inbox, but I'm a chicken.
func WithDebugChannel() ConvergerOpt {
	return func(converger *converger) {
		converger.debug = make(chan debugReq)
	}
}

// WithMaxInboxSize sets the maximum backlog for converge requests. Converge
// requests surpassing the backlog will return an error
func WithMaxInboxSize(size int) ConvergerOpt {
	return func(converger *converger) {
		converger.inbox = make(chan message, size)
	}
}

// StartConverger starts a converger. Options may be provided to override default
// values for the compiler and inbox size.
// NOTE: A Converger runs in a goroutine, and this function launches it
func StartConverger(opts ...ConvergerOpt) Converger {
	converger := &converger{
		inbox:        make(chan message, 5),
		stop:         make(chan struct{}, 1),
		currentState: &idle{},
		selfHuper:    &selfHuper{},
		compiler:     NewPhaseOrderedCompiler(),
	}
	for _, opt := range opts {
		opt(converger)
	}
	go converger.run()
	return converger
}

func (converger *converger) run() {
	for {
		select {
		case d := <-converger.debug:
			d.Do(converger)
		default:
		}

		select {
		case <-converger.stop:
			//TODO(jaym) drain inbox and fail pending requests
			close(converger.inbox)
			close(converger.stop)
			return
		case req := <-converger.inbox:
			logrus.Debugf("got message %s while %s", req, converger.currentState)
			nextState, stopTimer, err := converger.currentState.ProcessMessage(converger, req)
			if stopTimer {
				converger.currentState.StopTimer()
			}
			converger.setState(nextState)
			req.SendResponse(err)
		case d := <-converger.debug:
			d.Do(converger)
		case <-converger.currentState.TimeoutChan():
			logrus.Debugf("timed out while %s", converger.currentState)
			nextState, _, err := converger.currentState.ProcessMessage(
				converger,
				&timeout{})
			if err != nil {
				logrus.WithError(err).Error("error processing timeout event")
			}
			converger.setState(nextState)
		}
	}
}

func (converger *converger) setState(s state) {
	if s.String() != converger.currentState.String() {
		logrus.WithFields(logrus.Fields{
			"old_state": converger.currentState,
			"new_state": s,
		}).Debug("converger state transition")
	}
	converger.currentState = s
}

// state represents the possible states of our converge
// loop. Currently, we use a small state machine to represent the
// possible "waiting" states that we might end up in as the result of
// a converge request.
type state interface {
	// ProcessMessage takes the converger and the most recently
	// received message. The returned state is the next state of
	// the fsm. The returned error may be send to the sender of
	// the message. The retruned bool indicates whether we should
	// reset the current state timeout before moving to the next
	// state.
	ProcessMessage(*converger, message) (state, bool, error)
	// TimeoutChan returns a channel that should delivery a
	// message when this state has timed out. If a message is
	// received on the timeout channel a timeout message will be
	// sent to the current state.
	TimeoutChan() <-chan time.Time
	// StopTimer should correctly stop any timers or resources
	// associated with the TimeoutChan.
	StopTimer()
	// String is a convenience method used for logging. states should
	// have uniq string descriptors so we can more easily
	// conditionally log.
	String() string
}

// message represents the messages that can be sent to our
// converger. Messages move the converger from one state to another,
// possibly triggering some system change.
type message interface {
	// SendResponse delivers err back to the original requestor if
	// the original requestor expects a response.
	SendResponse(err error)
	// String is a convenience method used for logging.
	String() string
}

// TODO(ssd) 2019-04-21: Should we just panic here? It is a programming error.
func UnknownMessageError(m message, s state) error {
	return errors.Errorf("unknown message %s received while in state %s", m, s)
}

//
// Messages
//

// timeout is a message sent to the current state when the current
// state's timeout has fired. It is generated by the converger and is
// not sent by the user code so there is no caller to respond to.
type timeout struct{}

func (t *timeout) SendResponse(_ error) {}
func (t *timeout) String() string       { return "timeout" }

// convergeRequest is a user-sendable message that instructs the
// converger to run the converge loop.
type convergeRequest struct {
	eventSink    EventSink
	desiredState DesiredState
	doneChan     chan error
}

func (c *convergeRequest) SendResponse(err error) {
	c.doneChan <- err
	close(c.doneChan)
}

func (c *convergeRequest) String() string { return "convergeRequest" }

// stopRequest is a user-sendable message that instructs the
// converger to stop all services.
type stopRequest struct {
	eventSink    EventSink
	desiredState DesiredState
	doneChan     chan error
}

func (s *stopRequest) SendResponse(err error) {
	s.doneChan <- err
	close(s.doneChan)
}

func (s *stopRequest) String() string { return "stopRequest" }

// reconfigureDone is a user-sendable message that instructs the
// converger that we've completed a reconfiguration request.
type reconfigureDone struct {
	doneChan chan error
}

func (r *reconfigureDone) SendResponse(err error) {
	if r.doneChan != nil {
		r.doneChan <- err
		close(r.doneChan)
	}
}

func (r *reconfigureDone) String() string { return "reconfigureDone" }

// States
//
// We use the following state machine to handle the two "waiting"
// states that we need to track inside the converger.
//
//        +-----------+  convergeRequest           +-------------------+
//        |           |  err = ErrRestartPending   |                   |
//        |           +--------------------------->|                   |
//        |    idle   |                            | waitingForRestart |
//        |           |<---------------------------+                   |
//        |           |       timeout              |                   |
//        +-------+---+                            +-------------------+
//            ^   |
//            |   | convergeRequest
//   timeout 2|   | err = ErrSelfReconfigurePending
//            |   |
//            |   v
//         +--+--------------------+
//         |                       | timeout 1
//         |                       +---------+
//         | waitingForReconfigure |         |
//         |                       |<--------+
//         |                       |
//         +-----------------------+
//

// idle is the default state of the converger. In this state we
// are ready to accept any user-sendable requests. The idle state
// never times out.
type idle struct{}

func (i *idle) ProcessMessage(c *converger, msg message) (state, bool, error) {
	switch msg.(type) {
	case *convergeRequest:
		req := msg.(*convergeRequest)
		convergePlan, err := c.compiler.Compile(req.desiredState)
		if err != nil {
			logrus.WithError(err).Error("Failed to compile")
			return i, false, err
		}
		err = convergePlan.Execute(req.eventSink)
		return i.handleError(err)
	case *stopRequest:
		req := msg.(*stopRequest)
		convergePlan, err := c.compiler.Compile(req.desiredState)
		if err != nil {
			logrus.WithError(err).Error("Failed to compile")
			return i, false, err
		}
		err = convergePlan.Execute(req.eventSink)
		return i.handleError(err)
	case *reconfigureDone:
		logrus.Warn("received reconfigureDone when in idle state. We should have been in waitingForReconfigure")
		return i, false, nil
	case *timeout:
		logrus.Warn("received timeout when in idle state.")
		return i, false, nil
	default:
		return i, false, UnknownMessageError(msg, i)
	}
}

func (i *idle) handleError(err error) (state, bool, error) {
	switch err {
	case api.ErrRestartPending:
		return newWaitingForRestart(), false, err
	case api.ErrSelfReconfigurePending:
		return newWaitingForReconfigure(false), false, err
	case api.ErrSelfUpgradePending:
		// NOTE(ssd) 2019-04-19: A self-upgrade is
		// essentially just a restart from hab-sup.
		return newWaitingForRestart(), false, err
	default:
		return i, false, err
	}
}

func (i *idle) TimeoutChan() <-chan time.Time { return nil }
func (i *idle) StopTimer()                    {}
func (i *idle) String() string                { return "idle" }

// waitingForRestart is a converger state that we enter whenever we
// expect to be restarted by an external signal.
type waitingForRestart struct {
	timeoutable
}

func newWaitingForRestart() *waitingForRestart {
	return &waitingForRestart{
		timeoutable{timeout: waitingForRestartTimeout},
	}
}

func (w *waitingForRestart) ProcessMessage(_ *converger, msg message) (state, bool, error) {
	switch msg.(type) {
	case *timeout:
		logrus.Info("timed out waiting for restart, re-entering idle state")
		return &idle{}, true, nil
	case *stopRequest, *convergeRequest:
		logrus.Infof("ignoring message %v because we are waiting to be restarted", msg)
		return w, false, api.ErrRestartPending
	case *reconfigureDone:
		logrus.Warn("reconfigureDone received when waiting for restart!")
		return w, false, api.ErrRestartPending
	default:
		logrus.Infof("ignoring message %v because we are waiting to be restarted", msg)
		return w, false, UnknownMessageError(msg, w)
	}
}

func (w *waitingForRestart) String() string { return "waiting for restart" }

// waitingForReconfigure is a converer state that we enter whenever we
// expect to have received a reconfiguration request from the Habitat
// supervisor.
type waitingForReconfigure struct {
	timeoutable
	sentSelfHup bool
}

func newWaitingForReconfigure(sentSelfHup bool) *waitingForReconfigure {
	return &waitingForReconfigure{
		timeoutable: timeoutable{timeout: waitingForReconfigureTimeout},
		sentSelfHup: sentSelfHup,
	}
}

func (w *waitingForReconfigure) ProcessMessage(c *converger, msg message) (state, bool, error) {
	switch msg.(type) {
	case *timeout:
		if !w.sentSelfHup {
			logrus.Warnf("Expected HUP from habitat but none received after %s, sending HUP to self",
				waitingForReconfigureTimeout)
			c.selfHuper.Hup()
			return newWaitingForReconfigure(true), true, nil
		}
		logrus.Warnf("Attempt to self-HUP also failed to result in a reconfigure after %s, giving up",
			waitingForReconfigureTimeout)
		return &idle{}, true, nil
	case *reconfigureDone:
		return &idle{}, true, nil
	case *stopRequest, *convergeRequest:
		logrus.Infof("ignoring message %v because we are waiting to be reconfigured", msg)
		return w, false, api.ErrSelfReconfigurePending
	default:
		return w, false, UnknownMessageError(msg, w)
	}
}

func (w *waitingForReconfigure) String() string { return "waiting for reconfigure" }

// timeoutable is can be embedded in state's that want to time out
// after a duration.
type timeoutable struct {
	timeout time.Duration
	timer   *time.Timer
}

func (t *timeoutable) TimeoutChan() <-chan time.Time {
	if t.timer == nil {
		t.timer = time.NewTimer(t.timeout)
	}
	return t.timer.C
}

func (t *timeoutable) StopTimer() {
	if t.timer == nil {
		return
	}
	if !t.timer.Stop() {
		<-t.timer.C
	}
}

// Converger Message API Functions

// ReconfigureDone instructs the converger that another part of the
// system has finished a a reconfiguration request.
func (converger *converger) ReconfigureDone(t *Task) error {
	msg := &reconfigureDone{}
	if t != nil {
		msg.doneChan = t.C
	}
	return converger.sendMessage(t, msg)
}

// Converge requests the converger try to reach the specified desired state.
// This function returns an error if the backlog is full. Otherwise, a Task with
// a channel that will receive a message will be returned.
func (converger *converger) Converge(task *Task, desiredState DesiredState, eventSink EventSink) error {
	return converger.sendMessage(task, &convergeRequest{
		eventSink:    eventSink,
		desiredState: desiredState, //TODO(jaym) make a copy
		doneChan:     task.C,
	})
}

// StopServices requests that the converger tries to stop all
// non-deployment services on the given target.
func (converger *converger) StopServices(t *Task, tgt target.Target, eventSink EventSink) error {
	topology := make(Topology)
	topology[tgt] = []Service{
		{
			Name:          deploymentServiceName,
			ConvergeState: Skip(),
		},
	}
	desiredState := NewDesiredState(topology,
		NewSkipSupervisorState(),
		[]habpkg.HabPkg{}, // ignore-list is empty because cleanup mode is disabled
		depot.DisabledGC)

	return converger.sendMessage(t, &stopRequest{
		eventSink:    eventSink,
		desiredState: desiredState, //TODO(jaym) make a copy
		doneChan:     t.C,
	})
}

// Stop stops the converger
func (converger *converger) Stop() {
	converger.stop <- struct{}{}
}

func (converger *converger) sendMessage(t *Task, msg message) error {
	select {
	case converger.inbox <- msg:
		logrus.WithField("converger_inbox_len", len(converger.inbox)).Debugf("Sent %s request to converger", msg)
		return nil
	default:
		if t != nil {
			close(t.C)
		}
		return errors.New("Inbox full")
	}
}

// Debug
//
// Our converger has a debug channel where we can send debug requests
// to query the current state of the converger. This is used in the
// test. An alternative approach here would be to add a mutex.
//
// NOTE(ssd) 2019-04-23: This currently does not do a deep copy of the
// states, so callers should take care not to mutate it!
type debugReq interface {
	Do(*converger)
}

type debugResp struct {
	state    state
	inboxLen int
}

type debugGetReq struct {
	responseChan chan debugResp
}

func (d debugGetReq) Do(c *converger) {
	d.responseChan <- debugResp{
		state:    c.currentState,
		inboxLen: len(c.inbox),
	}
}

type debugSetReq struct {
	newState state
	doneChan chan struct{}
}

func (d debugSetReq) Do(c *converger) {
	logrus.WithFields(logrus.Fields{
		"new_state": d.newState,
	}).Debug("setting state via converger debug message")
	c.setState(d.newState)
	close(d.doneChan)
}

func (converger *converger) debugGetState() state {
	respChan := make(chan debugResp, 1)
	converger.debug <- debugGetReq{
		responseChan: respChan,
	}
	resp := <-respChan
	return resp.state
}

func (converger *converger) debugSetState(newState state) {
	done := make(chan struct{})
	converger.debug <- debugSetReq{
		doneChan: done,
		newState: newState,
	}
	<-done
}
