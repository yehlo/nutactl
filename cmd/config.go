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

import "github.com/spf13/cobra"

func newConfigCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Manage config",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runConfig),
	}
	cmd.AddCommand(
		newConfigContextCreateCommand(cli),
		newConfigContextSetCommand(cli),
		newConfigContextDeleteCommand(cli),
		newConfigContextListCommand(cli),
	)

	cmd.Flags().SortFlags = false
	return cmd
}

func runConfig(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}