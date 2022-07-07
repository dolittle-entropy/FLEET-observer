/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package config

import (
	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// SetupFor loads the koanf.Koanf configuration and creates a zerolog.Logger for the supplied
// cobra.Command
func SetupFor(cmd *cobra.Command) (*koanf.Koanf, zerolog.Logger, error) {
	config, err := LoadConfigFor(cmd)
	if err != nil {
		return nil, zerolog.Nop(), err
	}

	logger, err := CreateLoggerUsing(cmd, config)
	if err != nil {
		return nil, zerolog.Nop(), err
	}

	return config, logger, nil
}
