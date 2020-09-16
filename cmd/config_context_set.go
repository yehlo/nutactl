// Copyright © 2020 Joshua Leuenberger
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"strconv"
	"github.com/simonfuhrer/nutactl/pkg"
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigContextSetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set-context context",
		Short:                 "sets an existing context active",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runConfigContextSet),
	}

	return cmd
}

func runConfigContextSet(cli *CLI, cmd *cobra.Command, args []string) error {
	newContext := args[0]
	idInt, err := strconv.Atoi(newContext)
	id := uint32(idInt)
	if err != nil {
		return err
	}

	config.File = cfgFile
	config.SetContext(id)

	fmt.Println("Context " + newContext + " set!")
	return nil
}