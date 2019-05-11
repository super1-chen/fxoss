package cmd

import (

	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	DEBUG bool
)

var rootCmd = &cobra.Command{
	Use:"fxoss",
	Short: "FXDATA OSS command line tool",
	Long:"FXDATA OSS command line tool for get list of all cds, show cds detail and ssh login cds server",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println(args)
		fmt.Println("this is root func")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&DEBUG, "debug", "d", false, "Show debug info")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
