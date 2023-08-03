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
	"bufio"
	"fmt"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
)

func newLoginCommand(in io.Reader, out io.Writer) *cobra.Command {
	debug := false
	c := &cobra.Command{
		Use:   "login",
		Short: "Login into account and safe credentials into config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(in)
			fmt.Fprint(out, "Enter username/email: ")
			username, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			username = strings.TrimSpace(username)
			fmt.Print("Enter Password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			password := strings.TrimSpace(string(bytePassword))
			var logger *zap.Logger
			if debug {
				logger, _ = zap.NewDevelopment()
			} else {
				logger, _ = zap.NewProduction()
			}
			client := jdownloader.NewClient(username, password, logger.Sugar())
			err = client.Connect()
			if err != nil {
				return err
			}
			defer func(client jdownloader.JdClient) {
				err := client.Disconnect()
				if err != nil {
					fmt.Fprintf(out, "Failed to disconnectLater client: %v\n", err)
				}
			}(client)
			var cfg configData
			cfg.Mail = &username
			cfg.Password = &password

			return saveConfig(&cfg)
		},
	}
	addDebugFlag(c.Flags(), &debug)
	return c
}
