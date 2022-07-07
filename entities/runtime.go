/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import (
	"fmt"
	"time"
)

type RuntimeVersionUID string

var RuntimeVersionType = "RuntimeVersion"

type RuntimeVersion struct {
	UID  RuntimeVersionUID `bson:"_id" json:"uid"`
	Type string            `bson:"_type" json:"type"`

	Properties struct {
		Major      int       `bson:"major" json:"major"`
		Minor      int       `bson:"minor" json:"minor"`
		Patch      int       `bson:"patch" json:"patch"`
		Prerelease string    `bson:"prerelease" json:"prerelease,omitempty"`
		Released   time.Time `bson:"released" json:"-"`
	} `bson:"properties" json:"properties"`

	Links struct {
	} `bson:"links" json:"-"`
}

func NewRuntimeVersionUID(major, minor, patch int, prerelease string) RuntimeVersionUID {
	if prerelease == "" {
		return RuntimeVersionUID(fmt.Sprintf("%v.%v.%v", major, minor, patch))
	}
	return RuntimeVersionUID(fmt.Sprintf("%v.%v.%v-%v", major, minor, patch, prerelease))
}

func NewRuntimeVersion(major, minor, patch int, prerelease string, released time.Time) RuntimeVersion {
	version := RuntimeVersion{}
	version.UID = NewRuntimeVersionUID(major, minor, patch, prerelease)
	version.Type = RuntimeVersionType
	version.Properties.Major = major
	version.Properties.Minor = minor
	version.Properties.Patch = patch
	version.Properties.Prerelease = prerelease
	version.Properties.Released = released
	return version
}
