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
	"net"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	foreman "github.com/simonfuhrer/nutactl/pkg/foreman"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newForemanHostCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS]",
		Short:                 "create a host",
		Aliases:               []string{"cre", "c"},
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runForemanHostCreate),
	}
	flags := cmd.Flags()
	flags.SortFlags = false
	flags.String("location", "", "location name or id (required)")
	flags.String("organisation", "", "organisation name or id")
	flags.String("os", "", "os name or id (required)")
	flags.String("domain", "", "domain name or id (required)")
	flags.String("environment", "", "environment name or id")
	flags.String("hostgroup", "", "hostgroup name or id")
	flags.String("puppet-ca", "", "puppet-ca name or id")
	flags.String("puppet", "", "puppet name or id")
	flags.String("subnet", "", "subnet name or id")
	flags.String("ip", "", "ip address")
	flags.String("mac", "", "mac address")
	flags.String("comment", "", "host comment")
	flags.String("compute-resource", "", "compute-resource name or id")
	flags.String("compute-profile", "", "compute-profile name or id")
	flags.Bool("build", true, "host build mode")
	flags.Bool("start-host", false, "start host after creation")
	flags.String("provisionmethod", "build", "build or image")
	markFlagsRequired(cmd, "domain", "os", "location")
	return cmd
}

func runForemanHostCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name := args[0]
	s := spinner.New(spinner.CharSets[6], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Creating VM %s", name)
	s.Start()
	mac := viper.GetString("mac")
	domainIDOrName := viper.GetString("domain")
	osIDOrName := viper.GetString("os")
	locationIDOrName := viper.GetString("location")
	environmentIDOrName := viper.GetString("environment")
	hostgroupIDOrName := viper.GetString("hostgroup")
	puppetIDOrName := viper.GetString("puppet")
	puppetcaIDOrName := viper.GetString("puppet-ca")
	subnetIDOrName := viper.GetString("subnet")
	ip := viper.GetString("ip")
	comment := viper.GetString("comment")
	buildmode := viper.GetBool("build")
	provisionMethod := viper.GetString("provisionmethod")
	computeResourceIDOrName := viper.GetString("compute-resource")
	computeProfileIDOrName := viper.GetString("compute-profile")
	startHost := viper.GetBool("start-host")

	var computeResource *foreman.ComputeResource
	if len(subnetIDOrName) == 0 && len(ip) > 0 {
		return fmt.Errorf("ip provided without a subnet")
	}

	if provisionMethod == "image" {
		if len(computeResourceIDOrName) == 0 {
			return fmt.Errorf("missing compute-resource id or name")
		}
		if len(computeProfileIDOrName) == 0 {
			return fmt.Errorf("missing compute-profile id or name")
		}
	}

	s.Suffix = fmt.Sprintf(" Creating VM %s --> Get Domain %s", name, domainIDOrName)
	domain, err := cli.ForemanClient().GetDomain(cli.Context, domainIDOrName)
	if err != nil {
		return err
	}

	s.Suffix = fmt.Sprintf(" Creating VM %s --> Search existing host", name)
	searchfilter := fmt.Sprintf("name == %s.%s", name, domain.Name)
	hostresponse, err := cli.ForemanClient().SearchHost(cli.Context, searchfilter)
	if err != nil {
		return err
	}
	if len(hostresponse.Results) > 0 {
		return fmt.Errorf("host %s.%s with ID %d already exists", name, domain.Name, hostresponse.Results[0].ID)
	}

	s.Suffix = fmt.Sprintf(" Creating VM %s --> Get location %s", name, locationIDOrName)
	location, err := cli.ForemanClient().GetLocation(cli.Context, locationIDOrName)
	if err != nil {
		return err
	}

	s.Suffix = fmt.Sprintf(" Creating VM %s --> Get OS %s", name, osIDOrName)
	os, err := cli.ForemanClient().GetOperatingSystem(cli.Context, osIDOrName)
	if err != nil {
		return err
	}
	// Workaround to get all properties
	os, err = cli.ForemanClient().GetOperatingSystemByID(cli.Context, os.ID)
	if err != nil {
		return err
	}

	if len(os.Architectures) == 0 {
		return fmt.Errorf("missing os architectur")
	}

	if len(mac) == 0 {
		mac, err = generateMac()
		if err != nil {
			return err
		}
	}

	request := foreman.HostRequest{
		Host: foreman.NewHostData{
			ForemanObject: foreman.ForemanObject{
				Name:       name,
				LocationID: location.ID,
			},
			Mac:               mac,
			DomainID:          domain.ID,
			OperatingsystemID: os.ID,
			ArchitectureID:    os.Architectures[0].ID,
			Build:             buildmode,
			ProvisionMethod:   provisionMethod,
			Managed:           true,
		},
	}

	if len(computeResourceIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get ComputeResource %s", name, computeResourceIDOrName)
		computeResource, err = cli.foremanclient.GetComputeResource(cli.Context, computeResourceIDOrName)
		if err != nil {
			return err
		}
		request.Host.ComputeResourceID = computeResource.ID

		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get ComputeProfile %s", name, computeProfileIDOrName)
		computeProfile, err := cli.foremanclient.GetComputeProfile(cli.Context, computeProfileIDOrName)
		if err != nil {
			return err
		}
		computeProfile, err = cli.foremanclient.GetComputeProfileByID(cli.Context, computeProfile.ID)
		if err != nil {
			return err
		}
		request.Host.ComputeProfileID = computeProfile.ID

		images, err := cli.foremanclient.SearchOperatingSystemImages(cli.Context, os, fmt.Sprintf("compute_resource==%s", computeResource.Name))
		if err != nil {
			return err
		}
		if len(images.Results) == 0 {
			return fmt.Errorf("missing compute resource image for os %s", os.Name)
		}
		storageDomains, err := cli.foremanclient.GetComputeResourceStorageDomains(cli.Context, computeResource, "")
		if err != nil {
			return err
		}
		if len(storageDomains.Results) == 0 {
			return fmt.Errorf("missing compute resource storage domains. check computeresource")
		}
		targetSR := storageDomains.Results[0].UUID

		var indexAttr int
		found := false
		for index, attr := range computeProfile.ComputeAttributes {
			if attr.ComputeProfileID == computeProfile.ID {
				indexAttr = index
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing computeprofile in computeresource")
		}

		request.Host.ComputeAttributes = &foreman.ComputeAttributesXenHost{
			ImageUUID:   images.Results[0].UUID,
			Start:       startHost,
			TargetSR:    targetSR,
			ConfigDrive: "1",
			VCPUsMax:    fmt.Sprintf("%d", computeProfile.ComputeAttributes[indexAttr].VMAttrs.VCPUsMax),
			MemoryMax:   fmt.Sprintf("%d", computeProfile.ComputeAttributes[indexAttr].VMAttrs.MemoryMax),
			MemoryMin:   fmt.Sprintf("%d", computeProfile.ComputeAttributes[indexAttr].VMAttrs.MemoryMin),
		}
	}

	if len(environmentIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get Environment %s", name, environmentIDOrName)
		environment, err := cli.foremanclient.GetEnvironment(cli.Context, environmentIDOrName)
		if err != nil {
			return err
		}
		request.Host.EnvironmentID = environment.ID
	}

	if len(hostgroupIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get Hostgroup %s", name, hostgroupIDOrName)
		hostgroup, err := cli.foremanclient.GetHostgroup(cli.Context, hostgroupIDOrName)
		if err != nil {
			return err
		}
		request.Host.HostgroupID = hostgroup.ID
	}

	if len(puppetIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get SmartProxy Puppet %s", name, puppetIDOrName)
		puppet, err := cli.foremanclient.GetSmartProxy(cli.Context, puppetIDOrName)
		if err != nil {
			return err
		}
		request.Host.PuppetProxyID = puppet.ID
	}

	if len(puppetcaIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get SmartProxy Puppet-CA %s", name, puppetcaIDOrName)
		puppetca, err := cli.foremanclient.GetSmartProxy(cli.Context, puppetcaIDOrName)
		if err != nil {
			return err
		}
		request.Host.PuppetCAProxyID = puppetca.ID
	}

	if len(subnetIDOrName) > 0 {
		s.Suffix = fmt.Sprintf(" Creating VM %s --> Get Subnet %s", name, subnetIDOrName)
		subnet, err := cli.foremanclient.GetSubnet(cli.Context, subnetIDOrName)
		if err != nil {
			return err
		}
		request.Host.SubnetID = subnet.ID
		if len(ip) > 0 {
			ipaddress := net.ParseIP(ip)
			_, ipNetwork, err := net.ParseCIDR(subnet.NetworkAddress)
			if err != nil {
				return err
			}
			if !ipNetwork.Contains(ipaddress) {
				return fmt.Errorf("ip not in network %s", subnet.NetworkAddress)
			}
			request.Host.IP = ip
			request.Host.InterfacesAttributes = &foreman.InterfacesAttributes{
				Primary: &foreman.NetInterface{
					Primary:    true,
					Provision:  true,
					Type:       "interface",
					Identifier: "eth0",
				},
			}
			if len(computeResourceIDOrName) > 0 {
				networks, err := cli.foremanclient.GetComputeResourceAvailableNetworks(cli.Context, computeResource, "")
				if err != nil {
					return err
				}
				if len(networks.Results) == 0 {
					return fmt.Errorf("missing Compute Resource Networks. Please check your Compute Resource")
				}
				networkUUID := ""
				for _, network := range networks.Results {
					if strings.Contains(network.Name, subnet.NetworkAddress) {
						networkUUID = network.UUID
						break
					}
				}
				if len(networkUUID) == 0 {
					return fmt.Errorf("compute resource network with address %s not found", subnet.NetworkAddress)
				}
				request.Host.InterfacesAttributes.Primary.ComputeAttributes = &foreman.ComputeAttributesXenNetwork{
					NetworkUUID: networkUUID,
				}
			}
		}

	}

	if buildmode {
		rootPass, err := generatePassword(12)
		if err != nil {
			return err
		}
		request.Host.RootPass = rootPass
		if len(os.PartitionTables) == 0 {
			return fmt.Errorf("missing partition table")
		}
		if len(os.Media) == 0 {
			return fmt.Errorf("missing media")
		}
		if len(os.Media) > 0 {
			request.Host.MediumID = os.Media[0].ID
		}
		if len(os.PartitionTables) > 0 {
			request.Host.PartitionTableID = os.PartitionTables[0].ID
		}

	}

	if len(comment) > 0 {
		request.Host.Comment = comment
	}

	s.Suffix = fmt.Sprintf(" Creating VM %s --> Final create Host", name)
	host, err := cli.ForemanClient().CreateHost(cli.Context, &request)
	if err != nil {
		return err
	}
	s.Stop()
	fmt.Printf("Host %s with ID %d created\n", host.Name, host.ID)

	return nil
}
