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
	"github.com/spf13/cobra"
	"io"
)

func NewRootCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "jdcli",
		Short: "jDownloader CLI tool",
	}
	c.ResetFlags()
	c.AddCommand(newLoginCommand(in, out))
	c.AddCommand(newLinksCommand(out))
	c.AddCommand(newDownloadsCommand(out))
	c.AddCommand(newDeviceCommand(out))
	return c
}
