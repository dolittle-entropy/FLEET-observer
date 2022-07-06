package observing

import (
	"dolittle.io/fleet-observer/kubernetes"
	"errors"
	"fmt"
)

var (
	WrongKindReceived           = fmt.Errorf("%w: received wrong resource kind", kubernetes.IrrecoverableError)
	CouldNotParseRuntimeVersion = fmt.Errorf("%w: could not parse runtime version", kubernetes.IrrecoverableError)
	MissingConfiguration        = fmt.Errorf("%w: missing configuration while computing hash", kubernetes.IrrecoverableError)
	PodOwnerNotFound            = errors.New("could not find owner replicaset of pod")
)

func ReceivedWrongType(received any, expected string) error {
	return fmt.Errorf("%w: expected %s but got %T", WrongKindReceived, expected, received)
}

func FailedToParseRuntimeVersion(image string) error {
	return fmt.Errorf("%w: %v", CouldNotParseRuntimeVersion, image)
}

func CouldNotFindConfiguration(suffix string) error {
	return fmt.Errorf("%w: missing the suffix %v", MissingConfiguration, suffix)
}
