package entities

import (
	"fmt"
	"time"
)

type RuntimeVersion struct {
	Major      int       `bson:"major"`
	Minor      int       `bson:"minor"`
	Patch      int       `bson:"patch"`
	Prerelease string    `bson:"prerelease"`
	Released   time.Time `bson:"released"`
}

func (v RuntimeVersion) VersionString() string {
	if v.Prerelease == "" {
		return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%v.%v.%v-%v", v.Major, v.Minor, v.Patch, v.Prerelease)
}
