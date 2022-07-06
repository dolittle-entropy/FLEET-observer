package mongo

import (
	"errors"
)

var (
	NoDatabaseConfigured = errors.New("no MongoDB database name configured in connection string")
)
