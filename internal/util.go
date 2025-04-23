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
	"errors"
	"fmt"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/spf13/pflag"
	"io"
)

func addDeviceFlag(fs *pflag.FlagSet, target *string) {
	fs.StringVar(target, "device", *target, "Device name to use for this operation")
}

func addDebugFlag(fs *pflag.FlagSet, target *bool) {
	fs.BoolVar(target, "debug", *target, "Debugging flag")
}

func addJsonFlag(fs *pflag.FlagSet, target *bool) {
	fs.BoolVar(target, "json", *target, "JSON output flag")
}

func pickDevice(client jdownloader.JdClient) (string, error) {
	devs, err := client.ListDevices()
	if err != nil {
		return "", err
	}
	if len(*devs) == 0 {
		return "", errors.New("no device available")
	}
	a := *devs
	return a[0].Name, err
}

func doWithDevice(debug bool, devname string, out io.Writer, fn func(device jdownloader.Device) error) error {
	c, err := getClient(debug)
	if err != nil {
		return err
	}
	err = c.Connect()
	if err != nil {
		return err
	}
	defer clientCloser(c, out)
	if len(devname) == 0 {
		devname, err = pickDevice(c)
		if err != nil {
			return err
		}
	}
	dev, err := c.Device(devname)
	if err != nil {
		return err
	}
	return fn(dev)
}

func clientCloser(client jdownloader.JdClient, out io.Writer) {
	err := client.Disconnect()
	if err != nil {
		fmt.Fprintf(out, "Failed to disconnect client: %v\n", err)
	}
}

func formatSize(bytes *int64) string {
	if bytes == nil {
		return "N/A"
	}
	const unit = 1024
	if *bytes < unit {
		return fmt.Sprintf("%d B", *bytes)
	}
	div, exp := int64(unit), 0
	for n := *bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(*bytes)/float64(div), "KMGTPE"[exp])
}

func formatEta(seconds *int64) (res string) {
	if seconds == nil {
		return "N/A"
	}
	s := *seconds
	days, s := s/86400, s%86400
	hrs, s := s/3600, s%3600
	mins, s := s/60, s%60

	if days > 0 {
		res += fmt.Sprintf("%d days ", days)
	}
	res += fmt.Sprintf("%2.2d:%2.2d:%2.2d", hrs, mins, s)
	return
}

func formatSpeed(speed *float64) string {
	if speed == nil {
		return "N/A"
	}
	var size = int64(*speed)
	return fmt.Sprintf("%s/s", formatSize(&size))
}

func compressUrl(url string) string {
	if len(url) > 80 {
		return url[0:80]
	}
	return url
}
