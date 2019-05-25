// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"time"

	"github.com/spf13/cobra"

	"github.com/super1-chen/fxoss/conf"
	"github.com/super1-chen/fxoss/app"
	"github.com/super1-chen/fxoss/utils"
)

var (
	long *bool
)

// cdsListCmd represents the cdsList command
var cdsListCmd = &cobra.Command{
	Use:     "cds-list",
	Short:   "get cds list",
	Long:    `fxoss cds-list show all cds information`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runCDSList,
	Args:    cobra.MaximumNArgs(1),
	Example: "fxoss cds-list -l",
}

func init() {
	rootCmd.AddCommand(cdsListCmd)
	long = cdsListCmd.Flags().BoolP("long", "l", false, "show list information as  format")
}

func runCDSList(cmd *cobra.Command, args []string) {
	var option string
	now := time.Now().UTC()
	config := conf.NewConfig()
	app, err := app.NewOssServer(now, config, Debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}

	if len(args) == 1 {
		option = args[0]
	}
	err = app.ShowCDSList(option, *long)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}
}
