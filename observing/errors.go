package observing

import (
	"errors"
	"fmt"
)

var (
	WrongKindReceived           = errors.New("received wrong resource kind")
	CouldNotParseRuntimeVersion = errors.New("could not parse runtime version")
	PodOwnerNotFound            = errors.New("could not find owner replicaset of pod")
)

func ReceivedWrongType(received any, expected string) error {
	return fmt.Errorf("%w: expected %s but got %T", WrongKindReceived, expected, received)
}

func FailedToParseRuntimeVersion(image string) error {
	return fmt.Errorf("%w: %v", CouldNotParseRuntimeVersion, image)
}
