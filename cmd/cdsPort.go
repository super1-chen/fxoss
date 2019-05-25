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

// cdsPortCmd represents the cdsPort command
var cdsPortCmd = &cobra.Command{
	Use:   "cds-port",
	Short: "show cds port information",
	Long: `fxoss cds-port sn`,
	Args: requiredSN,
	PreRunE: func(cmd *cobra.Command, args []string) error {return  app.CheckEnvironment()},
	Run: runShowPort,
}


func init() {
	rootCmd.AddCommand(cdsPortCmd)
}

func runShowPort(cmd *cobra.Command, args []string){
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, Debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}

	err = app.ShowCDSPort(args[0])
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}
