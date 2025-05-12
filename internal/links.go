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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/spf13/cobra"
)

var (
	clCols = []string{"ID", "Name", "URL", "Status", "Size"}
)

func newLinksCommand(out io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "links",
		Short: "Interacts with links collector",
	}
	c.AddCommand(newAddLinksCommand(out))
	c.AddCommand(newListLinksCommand(out))
	return c
}

func newAddLinksCommand(out io.Writer) *cobra.Command {
	type addData struct {
		commonData
		fromFile    string
		links       []string
		packageName string
		downloadDir string
		autoStart   bool
	}
	var data addData
	data.autoStart = false
	data.links = make([]string, 0)
	c := &cobra.Command{
		Use:   "add",
		Short: "Add one or more links to LinkCollector",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(data.fromFile) > 0 {
				links, err := parseLinksFromFile(data.fromFile)
				if err != nil {
					return err
				}
				data.links = append(data.links, links...)
			}
			if len(data.links) == 0 {
				return errors.New("no links specified")
			}
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				opts := make([]jdownloader.AddLinksOptions, 0)
				opts = append(opts, jdownloader.AddLinksOptionAutostart(data.autoStart))
				if len(data.packageName) > 0 {
					opts = append(opts, jdownloader.AddLinksOptionPackage(data.packageName))
				}
				if len(data.downloadDir) > 0 {
					opts = append(opts, jdownloader.AddLinksOptionDestinationDir(data.downloadDir))
				}
				resp, err := dev.LinkGrabber().Add(data.links, opts...)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "Response: %v\n", resp.Data)
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	c.Flags().StringArrayVar(&data.links, "link", data.links, "Link to add. Can be specified multiple times")
	c.Flags().StringVar(&data.fromFile, "from-file", data.fromFile, "Path to file which contains URL on each line. Lines starting with ';' will be ignored")
	c.Flags().StringVar(&data.downloadDir, "download-dir", data.downloadDir, "Directory where to download files")
	c.Flags().StringVar(&data.packageName, "package-name", data.packageName, "Name of download package")
	c.Flags().BoolVar(&data.autoStart, "auto-start", data.autoStart, "Flag to determine whether files should start to download immediately or not")
	return c
}

func newListLinksCommand(out io.Writer) *cobra.Command {
	type listData struct {
		commonData
		json bool
	}
	var data listData
	c := &cobra.Command{
		Use:   "list",
		Short: "List all links in LinkGrabber",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				links, err := dev.LinkGrabber().Links()
				if err != nil {
					return err
				}

				if data.json {
					jsonObj, err := json.MarshalIndent(links, "", "    ")
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", jsonObj)
					return nil
				}

				if len(*links) > 0 {
					tbl := tablewriter.NewWriter(os.Stdout)
					tbl.Header(clCols)

					for _, link := range *links {
						row := make([]string, len(clCols))
						row[0] = strconv.FormatUint(uint64(*link.Uuid), 10)
						row[1] = *link.Name
						row[2] = compressUrl(*link.Url)
						if link.Status != nil {
							row[3] = *link.Status
						}
						if link.BytesTotal != nil {
							size := int64(*link.BytesTotal)
							row[4] = formatSize(&size)
						}
						tbl.Append(row)
					}
					return tbl.Render()
				} else {
					fmt.Fprintf(out, "No links\n")
				}
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	addJsonFlag(c.Flags(), &data.json)
	return c
}

func parseLinksFromFile(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		_, err = url.Parse(line)
		if err == nil {
			res = append(res, line)
		}
	}
	return res, nil
}
