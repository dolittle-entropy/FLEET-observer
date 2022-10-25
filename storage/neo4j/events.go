/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Events struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewEvents(session neo4j.SessionWithContext, ctx context.Context) *Events {
	return &Events{
		session: session,
		ctx:     ctx,
	}
}

func (e *Events) Set(event entities.Event) error {
	return runMultiUpdate(
		e.session,
		e.ctx,
		map[string]any{
			"uid":               event.UID,
			"count":             event.Properties.Count,
			"firstTime":         event.Properties.FirstTime.Format(time.RFC3339),
			"lastTime":          event.Properties.LastTime.Format(time.RFC3339),
			"platform":          event.Properties.Platform,
			"link_instance_uid": event.Links.HappenedToDeploymentInstanceUID,
		},
		`
			MERGE (event:`+event.Type+`:Event { _uid: $uid })
			SET event = { _uid: $uid, count: $count, firstTime: $firstTime, lastTime: $lastTime, platform: $platform }
			RETURN id(event)
		`,
		`
			MATCH (event:Event { _uid: $uid })
			WITH event
				MERGE (instance:DeploymentInstance { _uid: $link_instance_uid })
				WITH event, instance
					MERGE (event)-[:HappenedTo]->(instance)
					WITH event, instance
						MATCH (event)-[r:HappenedTo]->(other)
						WHERE other._uid <> instance._uid
						DELETE r
			RETURN id(event)
		`)
}

func (e *Events) Get(id entities.EventUID) (*entities.Event, bool, error) {
	event := &entities.Event{}
	found, err := findSingleJson(
		e.session,
		e.ctx,
		map[string]any{
			"uid": id,
		},
		`
			MATCH (event:Event { _uid: $uid })-[:HappenedTo]->(instance:DeploymentInstance)
			WITH {
				uid: event._uid,
				type: apoc.coll.removeAll(labels(event), ["Event"])[0],
				properties: {
					count: event.count,
					firstTime: toString(event.firstTime),
					lastTime: toString(event.lastTime),
					platform: event.platform
				},
				links: {
					happenedTo: instance._uid
				}
			} as entry
			RETURN apoc.convert.toJson(entry) as json
		`,
		event)
	return event, found, err
}

func (e *Events) List() ([]entities.Event, error) {
	var events []entities.Event
	return events, querySingleJson(
		e.session,
		e.ctx,
		`
			MATCH (event:Event)-[:HappenedTo]->(instance:DeploymentInstance)
			WITH {
				uid: event._uid,
				type: apoc.coll.removeAll(labels(event), ["Event"])[0],
				properties: {
					count: event.count,
					firstTime: toString(event.firstTime),
					lastTime: toString(event.lastTime),
					platform: event.platform
				},
				links: {
					happenedTo: instance._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&events)
}
