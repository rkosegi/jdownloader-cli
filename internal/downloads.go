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
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/spf13/cobra"
)

var (
	dlCols  = []string{"ID", "URL", "State", "ETA", "Speed", "Size"}
	pkgCols = []string{"ID", "Name", "Status", "Save to", "Total size"}
)

type commonData struct {
	debug  bool
	device string
}

func newDownloadsCommand(out io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "download",
		Short: "Manages downloads",
	}
	c.AddCommand(newDownloadLinksCommand(out))
	c.AddCommand(newDownloadPackageCommand(out))
	c.AddCommand(newDownloadStatusCommand(out))
	c.AddCommand(newDownloadCleanCommand(out))
	c.AddCommand(newDownloadPauseCommand(out))
	c.AddCommand(newDownloadStopCommand(out))
	c.AddCommand(newDownloadStartCommand(out))
	return c
}

func newDownloadLinksCommand(out io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "link",
		Short: "Manages download links",
	}
	c.AddCommand(newDownloadLinkListCommand(out))
	c.AddCommand(newDownloadLinkRmCommand(out))
	return c
}

func newDownloadLinkListCommand(out io.Writer) *cobra.Command {
	type newData struct {
		commonData
		json bool
	}
	var data newData
	c := &cobra.Command{
		Use:   "list",
		Short: "List downloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(device jdownloader.Device) error {
				links, err := device.Downloader().Links()
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

				tbl := tablewriter.NewWriter(os.Stdout)
				tbl.SetHeader(dlCols)

				for _, link := range *links {
					row := make([]string, len(dlCols))
					row[0] = strconv.FormatUint(uint64(*link.Uuid), 10)
					row[1] = compressUrl(*link.Url)
					if link.Status != nil {
						row[2] = *link.Status
					}
					row[3] = formatEta(link.Eta)
					row[4] = formatSpeed(link.Speed)
					row[5] = formatSize(link.BytesTotal)
					tbl.Append(row)
				}
				tbl.Render()
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	addJsonFlag(c.Flags(), &data.json)
	return c
}

func newDownloadLinkRmCommand(out io.Writer) *cobra.Command {
	type rmData struct {
		commonData
		id []int64
	}
	var data rmData
	data.id = make([]int64, 0)
	c := &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(data.id) == 0 {
				return errors.New("no link identifier(s) was specified (use --ids id1 --ids id2 ... )")
			}
			return doWithDevice(data.debug, data.device, out, func(device jdownloader.Device) error {
				return device.Downloader().Remove(data.id, []int64{})
			})
		},
	}
	c.Flags().Int64SliceVar(&data.id, "id", data.id, "Link identifier")
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}

func newDownloadPackageCommand(out io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "package",
		Short: "Manages download packages",
	}
	c.AddCommand(newDownloadPackageListCommand(out))
	return c
}

func newDownloadPackageListCommand(out io.Writer) *cobra.Command {
	type newData struct {
		commonData
		json bool
	}
	var data newData
	c := &cobra.Command{
		Use:   "list",
		Short: "List download packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(device jdownloader.Device) error {
				pkgs, err := device.Downloader().Packages()
				if err != nil {
					return err
				}

				if data.json {
					jsonObj, err := json.MarshalIndent(pkgs, "", "    ")
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", jsonObj)
					return nil
				}

				tbl := tablewriter.NewWriter(os.Stdout)
				tbl.SetHeader(pkgCols)

				for _, pkg := range *pkgs {
					row := make([]string, len(pkgCols))
					row[0] = strconv.FormatUint(uint64(*pkg.Uuid), 10)
					row[1] = *pkg.Name
					if pkg.Status != nil {
						row[2] = *pkg.Status
					}
					row[3] = *pkg.SaveTo
					row[4] = formatSize(pkg.BytesTotal)
					tbl.Append(row)
				}
				tbl.Render()
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	addJsonFlag(c.Flags(), &data.json)
	return c
}

func newDownloadStatusCommand(out io.Writer) *cobra.Command {
	var data commonData
	c := &cobra.Command{
		Use:   "status",
		Short: "Show downloader status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(device jdownloader.Device) error {
				si, err := device.Downloader().Speed()
				if err != nil {
					return err
				}
				st, err := device.Downloader().State()
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "Download status: %s\n", *st.State)
				fmt.Fprintf(out, "Download speed: %s\n", formatSpeed(si.Speed))
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}

func newDownloadCleanCommand(out io.Writer) *cobra.Command {
	var data commonData
	c := &cobra.Command{
		Use:   "clean",
		Short: "Clean completed downloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				dl := dev.Downloader()
				links, err := dl.Links()
				if err != nil {
					return err
				}
				toRemove := make([]int64, 0)
				for _, link := range *links {
					if link.Status != nil && *link.Status == "Finished" {
						fmt.Fprintf(out, "%s is completed and will be removed\n", *link.Url)
						toRemove = append(toRemove, *link.Uuid)
					}
				}
				if len(toRemove) > 0 {
					err = dl.Remove(toRemove, []int64{})
					if err != nil {
						return err
					} else {
						fmt.Fprintf(out, "%d links cleaned\n", len(toRemove))
					}
				} else {
					fmt.Fprintf(out, "Nothing to clean\n")
				}
				return nil
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}

func newDownloadPauseCommand(out io.Writer) *cobra.Command {
	var data commonData
	c := &cobra.Command{
		Use:   "pause",
		Short: "Pauses download",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				res, err := dev.Downloader().Pause()
				fmt.Fprintf(out, "Result : %t", res)
				return err
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}

func newDownloadStopCommand(out io.Writer) *cobra.Command {
	var data commonData
	c := &cobra.Command{
		Use:   "stop",
		Short: "Stops download",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				res, err := dev.Downloader().Stop()
				fmt.Fprintf(out, "Result : %t", res)
				return err
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}

func newDownloadStartCommand(out io.Writer) *cobra.Command {
	var data commonData
	c := &cobra.Command{
		Use:   "start",
		Short: "Starts a download",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doWithDevice(data.debug, data.device, out, func(dev jdownloader.Device) error {
				res, err := dev.Downloader().Stop()
				fmt.Fprintf(out, "Result : %t", res)
				return err
			})
		},
	}
	addDebugFlag(c.Flags(), &data.debug)
	addDeviceFlag(c.Flags(), &data.device)
	return c
}
