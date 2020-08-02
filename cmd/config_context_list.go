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
	"github.com/simonfuhrer/nutactl/cmd/displayers"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigContextListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list-context",
		Short:                 "sets an existing context active",
		Aliases:               []string{"l", "li"},
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runConfigContextList),
	}

	return cmd
}

func runConfigContextList(cli *CLI, cmd *cobra.Command, args []string) error {
	// get all IDs
	ids := make([]string, 0)
	all := viper.AllKeys()
	for _, value := range all {
		if strings.Contains(value, "id"){
			ids = append(ids, strings.Split(value, ".")[1])
		}
	}

	// create configContext structs
	configContexts := make(displayers.ConfigContextSlice, len(ids))
	for i, id := range ids {
		c := displayers.ConfigContext{
			ID:	id,
			URL: "url",
			User: "user",
			Insecure: "insecure",
		}
		
		configContexts[i] = c
	}

	return outputResponse(configContexts)
}