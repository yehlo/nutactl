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
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/simonfuhrer/nutactl/pkg"
)

func newConfigContextCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create-context",
		Short:                 "Creates a new context and activates it",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runConfigContextCreate),
	}
	return cmd
}

func getUserInput(message string) (text string, err error) {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	text, err = reader.ReadString('\n') // break on newline
	text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF

	if err != nil {
		return "", err
	}

	return strings.ToLower(text), nil
}

func runConfigContextCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	
	// query user for context url
	url, err := getUserInput("URL to use: ")
	if err != nil {
		return err
	}

	// query user for context user
	user, err := getUserInput("User to use: ")
	if err != nil {
		return err
	}

	// query user for context pass
	pw, err := readUserPW()
	if err != nil {
		return err
	}

	fmt.Println("")

	// query user for context security
	ntxContextInsecureStr, err := getUserInput("Accept insecure TLS certificates (y/N): ")
	if err != nil {
		return err
	}
	if ntxContextInsecureStr == "y" {
		ntxContextInsecureStr = "true"
	} else {
		ntxContextInsecureStr = "false"
	}

	// convert string to bool
	insecure, err := strconv.ParseBool(ntxContextInsecureStr)
	if err != nil {
		return err
	}

	fmt.Println("setting File")
	config.File = cfgFile
	fmt.Println("cfgFile")
	fmt.Println(cfgFile)
	id := int(config.CreateContext(url, user, pw, insecure))
	newContext := strconv.Itoa(id)
	fmt.Println("newContext: " + newContext)


	// write other data to configfile
	// configContextRoot := fmt.Sprintf("ntxContexts.%d", ntxContextID)
	// viper.Set(configContextRoot + ".id", ntxContextID)
	// viper.Set(configContextRoot + ".url", ntxContextURL)
	// viper.Set(configContextRoot + ".user", ntxContextUser)
	// viper.Set(configContextRoot + ".insecure", ntxContextInsecure)

	// fmt.Println(cfgFile)
	// fmt.Println(viper.ConfigFileUsed())

	// // save config to cfgfile
	// viper.WriteConfig()

	// // activate context (cli Client automatically reads in env variables)
	// viper.Set("ntxContexts.active", ntxContextID)
	// os.Setenv(appName + "_API_URL", ntxContextURL)
	// os.Setenv(appName + "_USERNAME", ntxContextUser)
	// os.Setenv(appName + "_PASSWORD", ntxContextPW)
	// os.Setenv(appName + "_INSECURE", ntxContextInsecureStr)

	return nil
}