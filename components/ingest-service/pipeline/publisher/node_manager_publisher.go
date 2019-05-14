package publisher

import (
	"github.com/chef/automate/components/ingest-service/pipeline/message"
	log "github.com/sirupsen/logrus"
)

func BuildNodeManagerPublisher(nodeManagerClient NodeManagerClient) message.ChefRunPipe {
	return func(in <-chan message.ChefRun) <-chan message.ChefRun {
		return nodeManagerPublisher(in, nodeManagerClient)
	}
}

func nodeManagerPublisher(in <-chan message.ChefRun, nodeManagerClient NodeManagerClient) <-chan message.ChefRun {
	maxNumberOfBundledMsgs := 100
	out := make(chan message.ChefRun, maxNumberOfBundledMsgs)
	go func() {
		for msg := range in {
			// send to node manager from here.
			log.Infof("publish to node manager %v", msg.Run)

			// nodeManagerClient.sendNode(msg.Run)

			out <- msg
		}
		close(out)
	}()

	return out
}
