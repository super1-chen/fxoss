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

	"github.com/super1-chen/fxoss/utils"
	"github.com/super1-chen/fxoss/conf"
	"github.com/super1-chen/fxoss/app"
)

// cdsShowCmd represents the cdsShow command
var cdsShowDetail = &cobra.Command{
	Use:   "cds-show",
	Short: "Show cds detail info",
	Long: `fxoss cds-show sn`,
	PreRunE: func(cmd *cobra.Command, args []string) error {return  app.CheckEnvironment()},
	Run: runShowDetail,
	Args: requiredSN,
}

func init() {
	rootCmd.AddCommand(cdsShowDetail)
}

func runShowDetail(cmd *cobra.Command, args []string){
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, Debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}
	err = app.ShowCDSDetail(args[0])
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}
