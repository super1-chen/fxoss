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
	"fmt"
	//github.com/spf13/cobra/cobra/cmd"

	"github.com/spf13/cobra"
)


var (
	long *bool

)

var longMessage = `Get and show all cds information from oss server

Example:
fxoss cds-list 南京大学
fxoss cds-list CAS50444471
`

// cdsListCmd represents the cdsList command
var cdsListCmd = &cobra.Command{
	Use:   "cds-list",
	Short: "Get cds list",
	Long: longMessage,
	Run: runCDSList,
	Args: cobra.MaximumNArgs(1),


}

func init() {
	rootCmd.AddCommand(cdsListCmd)
    long = cdsListCmd.Flags().BoolP("long", "l", false, "show list information as long format")
}


func runCDSList(cmd *cobra.Command, args []string){
	if *long == true {
		fmt.Println("cds-list long called")
	} else {
		fmt.Println("cds-list short called")
	}
}
