package server

import (
	"context"

	"github.com/sirupsen/logrus"

	dls "github.com/chef/automate/api/interservice/data_lifecycle"
	es "github.com/chef/automate/api/interservice/es_sidecar"
)

type DataLifecycleManageableServer struct {
	purgeClient   es.EsSidecarClient
	purgePolicies []PurgePolicy
}

type PurgePolicy struct {
	// PurgeOlderThanDays describes how long we keep data for either time series or non-time series indices. Anything
	// older than this is allowed to be regularly removed.
	PurgeOlderThanDays int32
	// IndexName is the base name of the time series index that should be purged
	IndexName string
	// Don't delete the entire index for non-time series databases; just delete documents over a given age from the index
	IsNonTimeSeries bool
	// IMPLEMENTED FOR NON-TIME SERIES INDICES ONLY!
	// Optional field. Purge documents based on the time value in this field rather than
	// the default field, end_time. Values for this field come from client purge policies.
	// If this field is empty, purge defaults to end_time.
	CustomPurgeField string
}

func NewDataLifecycleManageableServer(purgeClient es.EsSidecarClient, purgePolicies []PurgePolicy) *DataLifecycleManageableServer {
	return &DataLifecycleManageableServer{
		purgeClient:   purgeClient,
		purgePolicies: purgePolicies,
	}
}

func (server *DataLifecycleManageableServer) Purge(ctx context.Context, request *dls.PurgeRequest) (*dls.PurgeResponse, error) {
	purgeResponse := &dls.PurgeResponse{
		Id:              request.GetId(),
		ComponentStatus: make(map[string]*dls.PurgeStatus),
	}

	for _, policy := range server.purgePolicies {
		indexName := policy.IndexName

		logctx := logrus.WithFields(logrus.Fields{
			"purge_id":   request.GetId(),
			"index_name": indexName,
		})

		logctx.WithField("older_than_days", policy.PurgeOlderThanDays).Info("Requesting purge")

		var (
			resp *es.PurgeResponse
			err  error
		)
		if policy.IsNonTimeSeries {
			logctx.Debug("Calling purge documents from index by age on data lifecycle component...")
			resp, err = server.purgeClient.PurgeDocumentsFromIndexByAge(ctx, &es.PurgeRequest{
				Id:               request.GetId(),
				Index:            indexName,
				OlderThanDays:    policy.PurgeOlderThanDays,
				CustomPurgeField: policy.CustomPurgeField,
			})
		} else {
			resp, err = server.purgeClient.PurgeTimeSeriesIndicesByAge(ctx, &es.PurgeRequest{
				Id:            request.GetId(),
				Index:         indexName,
				OlderThanDays: policy.PurgeOlderThanDays,
			})
		}

		if err != nil {
			if policy.IsNonTimeSeries {
				logctx.WithError(err).Error("PurgeDocumentsFromIndexByAge rpc call failed")
			} else {
				logctx.WithError(err).Error("PurgeTimeSeries rpc call failed")
			}

			// Errors are transformed into a purge status code and a message. We don't return
			// a grpc error so that we can allow multiple components to have their own status
			purgeResponse.ComponentStatus[indexName] = &dls.PurgeStatus{
				Status: dls.PurgeStatus_FAILED,
				Msg:    err.Error(),
			}
		} else {

			if resp.GetSuccess() {
				logctx.Info("Purge Success")

				purgeResponse.ComponentStatus[indexName] = &dls.PurgeStatus{
					Status: dls.PurgeStatus_SUCCESS,
					Msg:    resp.GetMessage(),
				}
			} else {
				logctx.Error("Purge Failed")

				purgeResponse.ComponentStatus[indexName] = &dls.PurgeStatus{
					Status: dls.PurgeStatus_FAILED,
					Msg:    resp.GetMessage(),
				}
			}
		}
	}

	return purgeResponse, nil
}
