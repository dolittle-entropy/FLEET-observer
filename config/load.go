/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"
	"strings"
)

// LoadConfigFor loads the configuration using the given cobra.Command flags,
// and any supplied YAML files through '--config' arguments
func LoadConfigFor(cmd *cobra.Command) (*koanf.Koanf, error) {
	k := koanf.New(".")

	configFiles, _ := cmd.Flags().GetStringSlice("config")
	for _, configFile := range configFiles {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	if err := k.Load(env.Provider("", k.Delim(), environmentVariableNameToKey), nil); err != nil {
		return nil, err
	}

	if err := k.Load(posflag.Provider(cmd.Flags(), k.Delim(), k), nil); err != nil {
		return nil, err
	}

	return k, nil
}

func environmentVariableNameToKey(name string) string {
	return strings.Replace(strings.ToLower(name), "_", ".", -1)
}
