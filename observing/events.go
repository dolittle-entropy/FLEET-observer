/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"fmt"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	listersAppsV1 "k8s.io/client-go/listers/apps/v1"
	listersCoreV1 "k8s.io/client-go/listers/core/v1"
	"strings"
)

type EventsHandler struct {
	events      *mongo.Events
	pods        listersCoreV1.PodLister
	replicasets listersAppsV1.ReplicaSetLister
	logger      zerolog.Logger
}

func NewEventsHandler(events *mongo.Events, pods listersCoreV1.PodLister, replicasets listersAppsV1.ReplicaSetLister, logger zerolog.Logger) *EventsHandler {
	return &EventsHandler{
		events:      events,
		pods:        pods,
		replicasets: replicasets,
		logger:      logger,
	}
}

func (eh *EventsHandler) Handle(obj any, _deleted bool) error {
	event, ok := obj.(*coreV1.Event)
	if !ok {
		return ReceivedWrongType(obj, "Event")
	}

	logger := eh.logger.With().Str("namespace", event.GetNamespace()).Str("name", event.GetName()).Logger()
	if event.InvolvedObject.Kind != "Pod" {
		logger.Trace().Msg("Skipping event because it does not involve a pod")
		return nil
	}

	pod, err := eh.pods.Pods(event.InvolvedObject.Namespace).Get(event.InvolvedObject.Name)
	if err != nil && errors.IsNotFound(err) {
		logger.Trace().Err(err).Msg("Skipping event because the pod no longer exists")
		return nil
	} else if err != nil {
		return err
	}

	tenantID, applicationID, environmentName, _, ok := GetMicroserviceIdentifiers(pod.ObjectMeta)
	if !ok {
		logger.Trace().Msg("Skipping event because the pod is missing microservice identifiers")
		return nil
	}

	_, _, ok = getRuntimeAndHeadContainer(pod.Spec)
	if !ok {
		logger.Trace().Msg("Skipping event because the pod does not have a runtime and head container")
		return nil
	}

	replicaSet, err := GetPodOwner(pod, eh.replicasets)
	if err != nil {
		logger.Trace().Msg("Skipping event because the pod owner could not be found")
		return nil
	}

	instanceUID := entities.NewDeploymentInstanceUID(
		tenantID,
		applicationID,
		environmentName,
		fmt.Sprintf("%v", replicaSet.GetGeneration()),
		string(pod.GetUID()),
	)

	return eh.handleDeploymentInstanceEvent(instanceUID, event, logger)
}

func (eh *EventsHandler) handleDeploymentInstanceEvent(id entities.DeploymentInstanceUID, event *coreV1.Event, logger zerolog.Logger) error {
	switch event.Reason {
	case "BackOff":
		return eh.handleBackOffEvent(id, event, logger)
	case "Created":
		return nil
	case "Failed":
		return nil
	case "Killed":
		return nil
	case "Pulled":
		return nil
	case "Pulling":
		return nil
	case "RELOAD":
		return nil
	case "Scheduled":
		return nil
	case "Started":
		return nil
	default:
		logger.Warn().Str("reason", event.Reason).Msg("Skipping event with unhandled reason")
		return nil
	}
}

func (eh *EventsHandler) handleBackOffEvent(id entities.DeploymentInstanceUID, event *coreV1.Event, logger zerolog.Logger) error {
	platformContainer := strings.Contains(event.InvolvedObject.FieldPath, "runtime")

	if strings.Contains(event.Message, "restarting") {
		event := entities.NewFailedToStartEvent(
			string(event.GetUID()),
			int(event.Count),
			event.FirstTimestamp.UTC(),
			event.LastTimestamp.UTC(),
			platformContainer,
			id,
		)
		if err := eh.events.Set(event); err != nil {
			return err
		}
		logger.Debug().Interface("event", event).Msg("Updated event")
		return nil
	}

	if strings.Contains(event.Message, "pulling") {
		event := entities.NewFailedToPullEvent(
			string(event.GetUID()),
			int(event.Count),
			event.FirstTimestamp.UTC(),
			event.LastTimestamp.UTC(),
			platformContainer,
			id,
		)
		if err := eh.events.Set(event); err != nil {
			return err
		}
		logger.Debug().Interface("event", event).Msg("Updated event")
		return nil
	}

	logger.Warn().Str("eventMessage", event.Message).Msg("Skipping BackOff event with unhandled message")
	return nil
}
