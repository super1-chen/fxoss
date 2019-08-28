package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/super1-chen/fxoss/app"
	"github.com/super1-chen/fxoss/conf"
	"github.com/super1-chen/fxoss/utils"
	"github.com/super1-chen/fxoss/version"
)

var (
	// global flag
	debug *bool
	// cds list partion
	long *bool
	// cds login partion
	r    *int
	frpc *bool
)

var rootCmd = &cobra.Command{
	Use:   "fxoss",
	Short: "fxoss is a command line tool for fxdata Ops team",
	Long:  "fxoss is a command line tool for get cds list, show cds detail and ssh login cds server...",
}

func init() {
	// show version
	rootCmd.AddCommand(versionCmd)
	debug = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "run fxoss in verbose mode")
	//
	rootCmd.AddCommand(cdsListCmd)
	long = cdsListCmd.Flags().BoolP("long", "l", false, "show list information as  format")
	// cds login partion
	rootCmd.AddCommand(cdsLoginCmd)
	frpc = cdsLoginCmd.Flags().BoolP("frpc", "F", false, "login cds in frpc mode")
	r = cdsLoginCmd.Flags().IntP("retry", "r", 3, "retry times of SSH login")
	// cds port partion
	rootCmd.AddCommand(cdsPortCmd)
	// show csd detail partion
	rootCmd.AddCommand(cdsShowDetail)
	// make cds report partion
	rootCmd.AddCommand(cdsReportShow)
	// make web root partion
	rootCmd.AddCommand(cdsWebRoot)
}

func requiredSN(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("cds sn is required")
	}
	return nil
}

func requiredValidEmail(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("one email address is required")
	}
	for _, arg := range args {
		if !(strings.HasSuffix(arg, "@fxdata.cn") || strings.HasPrefix(arg, "@ifeixiang.com")) {
			return fmt.Errorf("illegal format email %s", arg)
		}
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
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

// cds list partion
var cdsListCmd = &cobra.Command{
	Use:     "cds-list",
	Short:   "Show cds list",
	Long:    `fxoss cds-list show all cds information`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runCDSList,
	Args:    cobra.MaximumNArgs(1),
	Example: "fxoss cds-list -l",
}

func runCDSList(cmd *cobra.Command, args []string) {
	var option string
	now := time.Now().UTC()
	config := conf.NewConfig()
	app, err := app.NewOssServer(now, config, *debug)
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

// cds login partion
var cdsLoginCmd = &cobra.Command{
	Use:     "cds-login",
	Short:   "SSH login remote server",
	Long:    `fxoss cds-login sn`,
	Args:    requiredSN,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runLoginCDS,
}

func runLoginCDS(cmd *cobra.Command, args []string) {
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, *debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}

	err = app.LoginCDS(args[0], *r, *frpc)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}

// cds port partion
var cdsPortCmd = &cobra.Command{
	Use:     "cds-port",
	Short:   "Show cds port information",
	Long:    `fxoss cds-port sn`,
	Args:    requiredSN,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runShowPort,
}

func runShowPort(cmd *cobra.Command, args []string) {
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, *debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}

	err = app.ShowCDSPort(args[0])
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}

// show cds detail partion
var cdsShowDetail = &cobra.Command{
	Use:     "cds-show",
	Short:   "Show cds detail info",
	Long:    `fxoss cds-show sn`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runShowDetail,
	Args:    requiredSN,
}

func runShowDetail(cmd *cobra.Command, args []string) {
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, *debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}
	err = app.ShowCDSDetail(args[0])
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}

// cdsShowCmd represents the cdsShow command
var cdsReportShow = &cobra.Command{
	Use:     "cds-report",
	Short:   "Make cds disk type report and send the report by email",
	Long:    `fxoss cds-report chenc@fxdata.cn chenc@ifeixiang.com`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runReport,
	Args:    requiredValidEmail,
}

func runReport(cmd *cobra.Command, args []string) {
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, *debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}
	err = app.ReportCDS(now, args...)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}
}

// cdsShowCmd represents the cdsShow command
var cdsWebRoot = &cobra.Command{
	Use:     "web-root",
	Short:   "Get cds web root password",
	Long:    `fxoss web-root sn`,
	PreRunE: func(cmd *cobra.Command, args []string) error { return app.CheckEnvironment() },
	Run:     runWebRoot,
	Args:    requiredSN,
}

func runWebRoot(cmd *cobra.Command, args []string) {
	now := time.Now().UTC()
	config := conf.NewConfig()

	app, err := app.NewOssServer(now, config, *debug)
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
		return
	}
	err = app.WebRoot(args[0])
	if err != nil {
		utils.ErrorPrintln(err.Error(), false)
	}

}
