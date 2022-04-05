package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/web"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	RootCmd = &cobra.Command{
		Use:  "cloud189",
		Long: "cloud189 enables users to manage cloud files through the command line. For more information, please visit https://github.com/gowsp/cloud189",
	}
)

func AddCommand(cmds ...*cobra.Command) {
	RootCmd.AddCommand(cmds...)
}
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cloud189/config.json)")

	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(signCmd)
	RootCmd.AddCommand(upCmd)
	RootCmd.AddCommand(rmCmd)
	RootCmd.AddCommand(dlCmd)
	RootCmd.AddCommand(lsCmd)
	RootCmd.AddCommand(mkdirCmd)
	RootCmd.AddCommand(mvCmd)
	RootCmd.AddCommand(cpCmd)
	RootCmd.AddCommand(dfCmd)
	RootCmd.AddCommand(webdavCmd)
}

var app pkg.App
var once sync.Once

func App() pkg.App {
	once.Do(func() {
		// web.Api
		api := web.NewApi(cfgFile)
		app = drive.NewClient(api)
	})
	return app
}
