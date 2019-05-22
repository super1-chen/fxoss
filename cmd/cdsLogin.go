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
	"github.com/super1-chen/fxoss/fxoss"
	"github.com/super1-chen/fxoss/utils"
)

var (
	r *int
	frpc *bool
)

// cdsLoginCmd represents the cdsLogin command
var cdsLoginCmd = &cobra.Command{
	Use:   "cds-login",
	Short: "ssh login remote server",
	Long: `fxoss cds-login sn`,
	Args: requiredSN,
	Run: runLoginCDS,
}


func init() {
	rootCmd.AddCommand(cdsLoginCmd)
	frpc = cdsLoginCmd.Flags().BoolP("frpc", "F", false, "login cds in frpc mode")
	r = cdsLoginCmd.Flags().IntP("retry", "r", 3, "retry times of SSH login")
}

func runLoginCDS(cmd *cobra.Command, args []string){
	now := time.Now().UTC()
	app, err := fxoss.NewOssServer(now, Debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}

	err = app.LoginCDS(args[0], *r, *frpc)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}
