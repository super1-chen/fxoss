package cmd

import (

	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/super1-chen/fxoss/version"
)

var (
	Debug bool
)

var rootCmd = &cobra.Command{
	Use:"fxoss",
	Short: "fxoss is a command line tool for fxdata Ops team",
	Long:"fxoss is a command line tool for get cds list, show cds detail and ssh login cds server...",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&Debug, "verbose", "v", false, "run fxoss in verbose mode")
}

func requiredSN(cmd *cobra.Command, args[]string) error {
	if len(args) != 1 {
		return fmt.Errorf("cds sn is required")
	}
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fxoss",
	Long:  `All software has versions. This is fxoss's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fxoss version is:", version.Version)
	},
}

// Execute run the command tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
