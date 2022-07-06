package observing

import (
	"errors"
	"fmt"
)

var (
	WrongKindReceived = errors.New("received wrong resource kind")
)

func ReceivedWrongType(received any, expected string) error {
	return fmt.Errorf("%w: expected %s but got %T", WrongKindReceived, expected, received)
}
