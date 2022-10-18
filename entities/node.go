/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import "fmt"

type NodeUID string

var NodeType = "Node"

type Node struct {
	UID  NodeUID `bson:"_id" json:"uid""`
	Type string  `bson:"_type" json:"type"`

	Properties struct {
		Hostname string `bson:"hostname" json:"hostname"`
		Image    string `bson:"image" json:"image"`
		Type     string `bson:"type" json:"type"`
	} `bson:"properties" json:"properties"`

	Link struct {
	} `bson:"links" json:"-"`
}

func NewNodeUID(nodename string) NodeUID {
	return NodeUID(fmt.Sprintf("%v", nodename))
}

func NewNode(nodename, hostname, image, nodetype string) Node {
	node := Node{}
	node.UID = NewNodeUID(nodename)
	node.Type = NodeType
	node.Properties.Hostname = hostname
	node.Properties.Image = image
	node.Properties.Type = nodetype
	return node
}
