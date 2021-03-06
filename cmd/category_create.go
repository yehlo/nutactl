// Copyright © 2020 Simon Fuhrer
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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tecbiz-ch/nutanix-go-sdk/schema"
)

func newCategoryCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS] CATEGORY",
		Short:                 "Create an category",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runCategoryCreate),
	}
	cmd.Flags().StringP("description", "d", "", "Description")

	return cmd
}

func runCategoryCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name := args[0]
	description := viper.GetString("description")

	req := &schema.CategoryKey{
		Name: name,
	}
	if description != "" {
		req.Description = description
	}
	result, err := cli.Client().Category.Create(cli.Context, req)
	if err != nil {
		return err
	}
	fmt.Printf("Category %s created\n", result.Name)

	return nil
}
