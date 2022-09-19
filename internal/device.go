/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var devCols = []string{"ID", "Type", "Name", "Status"}

func newDeviceCommand(out io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "device",
		Short: "Manages devices",
	}
	c.AddCommand(newDeviceListCommand(out))
	return c
}

func newDeviceListCommand(out io.Writer) *cobra.Command {
	type listData struct {
		debug bool
	}
	var data listData
	c := &cobra.Command{
		Use:   "list",
		Short: "List all devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(data.debug)
			if err != nil {
				return err
			}
			err = c.Connect()
			if err != nil {
				return err
			}
			defer func(client jdownloader.JdClient) {
				err := client.Disconnect()
				if err != nil {
					fmt.Fprintf(out, "Failed to disconnect client: %v\n", err)
				}
			}(c)

			devs, err := c.ListDevices()
			if err != nil {
				return err
			}

			tbl := tablewriter.NewWriter(os.Stdout)
			tbl.SetHeader(devCols)
			for _, dev := range *devs {
				row := make([]string, len(devCols))
				row[0] = dev.Id
				row[1] = dev.Type
				row[2] = dev.Name
				row[3] = dev.Status
				tbl.Append(row)
			}
			tbl.Render()
			return nil
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	return c
}
