/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package mongo

import (
	"errors"
)

var (
	NoDatabaseConfigured = errors.New("no MongoDB database name configured in connection string")
)
