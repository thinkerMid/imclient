/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"github.com/spf13/cobra"
	"labs/src/sms"
	"labs/src/wsa/client"
	"os"
	"os/signal"
	"syscall"
)

// wsaagentCmd represents the wsaagent command
var wsaagentCmd = &cobra.Command{
	Use:   "wsaagent",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := client.CreateClient()

		go func() {
			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
			<-signalChannel

			cli.Cancel()
		}()

		sms.Instance().QueryCountries()

		cli.RunClient("IN")
		//fmt.Println((&WSAVersionCheck{}).checkVersionUpgrade())
	},
}

func init() {
	rootCmd.AddCommand(wsaagentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsaagentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsaagentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
