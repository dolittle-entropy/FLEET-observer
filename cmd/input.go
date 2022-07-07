/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cmd

import (
	"bufio"
	"context"
	"errors"
	"github.com/spf13/cobra"
)

var (
	ErrReadAborted = errors.New("reading input aborted")
	ErrNoInput     = errors.New("no input available")
)

func ReadLineFromInput(cmd *cobra.Command, ctx context.Context) (string, error) {
	readCh := make(chan string)
	errCh := make(chan error)

	go func() {
		scanner := bufio.NewScanner(cmd.InOrStdin())
		if !scanner.Scan() {
			errCh <- ErrNoInput
			return
		}
		readCh <- scanner.Text()
	}()

	select {
	case <-ctx.Done():
		return "", ErrReadAborted
	case err := <-errCh:
		return "", err
	case str := <-readCh:
		return str, nil
	}
}
