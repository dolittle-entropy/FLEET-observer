/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import (
	"fmt"
	"time"
)

type EventUID string

type Event struct {
	UID  EventUID `bson:"_id" json:"uid"`
	Type string   `bson:"_type" json:"type"`

	Properties struct {
		Count     int       `bson:"count" json:"count"`
		FirstTime time.Time `bson:"first_time" json:"firstTime"`
		LastTime  time.Time `bson:"last_time" json:"lastTime"`
		Platform  bool      `bson:"platform" json:"platform"`
	} `bson:"properties" json:"properties"`

	Links struct {
		HappenedToDeploymentInstanceUID DeploymentInstanceUID `bson:"happened_to_deployment_instance_uid" json:"happenedTo"`
	}
}

func newEvent(id EventUID, eventType string, count int, firstTime, lastTime time.Time, platform bool, instance DeploymentInstanceUID) Event {
	event := Event{}
	event.UID = id
	event.Type = eventType
	event.Properties.Count = count
	event.Properties.FirstTime = firstTime
	event.Properties.LastTime = lastTime
	event.Properties.Platform = platform
	event.Links.HappenedToDeploymentInstanceUID = instance
	return event
}

func NewKubernetesEventUID(eventID string) EventUID {
	return EventUID(fmt.Sprintf("kubernetes/%v", eventID))
}

var FailedToStartEventType = "FailedToStartEvent"

func NewFailedToStartEvent(eventID string, count int, firstTime, lastTime time.Time, platform bool, instance DeploymentInstanceUID) Event {
	return newEvent(
		NewKubernetesEventUID(eventID),
		FailedToStartEventType,
		count,
		firstTime,
		lastTime,
		platform,
		instance,
	)
}

var FailedToPullEventType = "FailedToPullEvent"

func NewFailedToPullEvent(eventID string, count int, firstTime, lastTime time.Time, platform bool, instance DeploymentInstanceUID) Event {
	return newEvent(
		NewKubernetesEventUID(eventID),
		FailedToPullEventType,
		count,
		firstTime,
		lastTime,
		platform,
		instance,
	)
}
