/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
)

type NodesHandler struct {
	nodes  *mongo.Nodes
	logger zerolog.Logger
}

func NewNodesHandler(nodes *mongo.Nodes, logger zerolog.Logger) *NodesHandler {
	return &NodesHandler{
		nodes:  nodes,
		logger: logger.With().Str("handler", "nodes").Logger(),
	}
}

func (nh *NodesHandler) Handle(obj any, _deleted bool) error {
	knode, ok := obj.(*coreV1.Node)
	if !ok {
		return ReceivedWrongType(obj, "Node")
	}

	logger := nh.logger.With().Str("node", knode.GetName()).Logger()

	hostname, ok := knode.GetLabels()["kubernetes.io/hostname"]
	if !ok {
		logger.Trace().Msg("Skipping node because it does not have an hostname annotation")
		return nil
	}

	image, ok := knode.GetLabels()["kubernetes.azure.com/node-image-version"]
	if !ok {
		logger.Trace().Msg("Skipping node because it does not have an node-image-version annotation")
		return nil
	}

	nodetype, ok := knode.GetLabels()["node.kubernetes.io/instance-type"]
	if !ok {
		logger.Trace().Msg("Skipping node because it does not have an instance-type annotation")
		return nil
	}

	node := entities.NewNode(knode.GetName(), hostname, image, nodetype)
	if err := nh.nodes.Set(node); err != nil {
		return err
	}
	logger.Debug().Interface("node", node).Msg("Updated node")

	return nil
}
