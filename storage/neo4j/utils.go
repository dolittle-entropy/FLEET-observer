/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	ErrResultRecordDoesNotContainJson = errors.New("the resulting record did not contain a field called 'json'")
	ErrResultJsonFieldWasNotString    = errors.New("the field called 'json' in the resulting record was not a string")
)

func runUpdate(session neo4j.SessionWithContext, ctx context.Context, cypher string, params map[string]any) error {
	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return err
	}

	_, err = result.Consume(ctx)
	return err
}

func runMultiUpdate(session neo4j.SessionWithContext, ctx context.Context, params map[string]any, cyphers ...string) error {
	_, err := session.ExecuteWrite(
		ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			for _, cypher := range cyphers {
				result, err := transaction.Run(ctx, cypher, params)
				if err != nil {
					return nil, err
				}

				_, err = result.Consume(ctx)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		})
	return err
}

func querySingleJson(session neo4j.SessionWithContext, ctx context.Context, cypher string, v any) error {
	result, err := session.Run(ctx, cypher, nil)
	if err != nil {
		return err
	}

	record, err := result.Single(ctx)
	if err != nil {
		return err
	}

	data, exists := record.Get("json")
	if !exists {
		return ErrResultRecordDoesNotContainJson
	}

	binary, ok := data.(string)
	if !ok {
		return ErrResultJsonFieldWasNotString
	}

	return json.Unmarshal([]byte(binary), v)
}
