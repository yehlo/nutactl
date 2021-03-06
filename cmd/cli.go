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
	"context"
	"fmt"
	"time"
	// "os"
	"github.com/briandowns/spinner"
	logrus "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	nutanix "github.com/tecbiz-ch/nutanix-go-sdk"
)

// PollIntervalinSeconds ..
const (
	PollIntervalinSeconds = 2
	appName               = "NUTACTL"
)

// CLI ...
type CLI struct {
	Endpoint            string
	Context             context.Context
	RootCommand         *cobra.Command
	client              *nutanix.Client
	millisecondsPerPoll time.Duration
	clusters            map[string]string
}

//NewCLI sadasdsa
func NewCLI() *CLI {
	// initConfig()
	cli := &CLI{
		Context:             context.Background(),
		millisecondsPerPoll: 1000 * PollIntervalinSeconds,
	}
	cli.RootCommand = NewRootCommand(cli)

	return cli
}

// wrapper func to bind all flags with viper and ensure a logout is perfomed
func (c *CLI) wrap(f func(*CLI, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		BindAllFlags(cmd)
		err := f(c, cmd, args)
		return checkErr(err)
	}
}

//Client cli ...
func (c *CLI) Client() *nutanix.Client {
	if c.client == nil {
		// check if api-url is set
		if viper.Get("api-url") == ""{
			panic("api-url not set!")
		}

		// pw is not managed by config file
		// if not set query user and define viper variable
		if ! viper.IsSet("password"){
			pw, err := readUserPW()
			if err != nil {
				panic(err)
			}
			viper.Set("password", pw)
			fmt.Println("")
		}

		configCreds := nutanix.Credentials{
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
		}
		opts := []nutanix.ClientOption{
			nutanix.WithCredentials(&configCreds),
			nutanix.WithEndpoint(viper.GetString("api-url")),
		}
		if viper.GetBool("insecure") {
			opts = append(opts, nutanix.WithInsecure())
		}

		logrus.Debugf("creating Nutanix Client")
		c.client = nutanix.NewClient(opts...)
	}
	return c.client
}

// InitAllClusters ...
func (c *CLI) InitAllClusters() error {
	logrus.Debugf("init Nutanix Clusters")
	if c.clusters == nil {
		clusters, err := c.client.Cluster.All(context.Background())
		if err != nil {
			return err
		}
		data := make(map[string]string)
		for _, cluster := range clusters.Entities {
			data[cluster.Metadata.UUID] = cluster.Spec.Name
		}
		c.clusters = data
	}
	return nil
}

//WaitTask ...
func (c *CLI) WaitTask(ctx context.Context, taskUUID string, timeoutSeconds int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(c.millisecondsPerPoll * time.Millisecond)
	defer ticker.Stop()
	s := spinner.New(spinner.CharSets[6], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Waiting Task ID %s", taskUUID)
	s.Start()
	for {
		select {
		case <-ticker.C:
			task, err := c.Client().Task.GetByUUID(ctx, taskUUID)
			if err != nil {
				return err
			}
			logrus.Debugf(*task.Status)
			switch *task.Status {
			case "SUCCEEDED":
				s.Stop()
				return nil
			case "FAILED":
				s.Stop()
				return fmt.Errorf(*task.ErrorDetail)

			}
		case <-ctx.Done():
			s.Stop()
			return fmt.Errorf("error waiting for task to be completed: %s", ctx.Err())
		}
	}
}
