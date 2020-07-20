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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigContextDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete-context",
		Short:                 "deletes an existing context active",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runConfigContextDelete),
	}

	return cmd
}

func runConfigContextDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	contextToDelete := args[0]

	// get user OK
	if askForConfirm(fmt.Sprintf("Delete %s ?", contextToDelete)) == nil {
		// check if it is currently used context, set active to nil if so 
		if viper.GetString("ntxContexts.active") == contextToDelete {
			fmt.Println("Currently used context is set to be deleted, no context is set as active")
			viper.Set("ntxContexts.active", "nil")

			// clear env
			os.Unsetenv(appName + "_API_URL")
			os.Unsetenv(appName + "_USERNAME")
			os.Unsetenv(appName + "_PASSWORD")
			os.Unsetenv(appName + "_INSECURE")
		}

		// remove context from viper cfgfile
		// viper is currently unable to use something like unset, setting whole context to nil as string
		viper.Set("ntxContexts." + contextToDelete, "nil")
		viper.WriteConfig()
		
		fmt.Println("context " + contextToDelete + " deleted!")
		return nil
	}
	return fmt.Errorf("operation aborted")
}