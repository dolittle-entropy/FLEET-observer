/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	return nil
}

func (e *Events) Get(id entities.EventUID) (*entities.Event, bool, error) {
	return nil, false, nil
}

func (e *Events) List() ([]entities.Event, error) {
	return nil, nil
}
