package cmd

import (
	"os"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/app"
	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/invoker"
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
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/cloud189/config.json)")

	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logoutCmd)
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
	RootCmd.AddCommand(shareCmd)
}

var singleton pkg.Drive
var once sync.Once

func App() pkg.Drive {
	once.Do(func() {
		if cfgFile == "" {
			cfgFile = invoker.DefaultPath()
		}
		api := app.New(cfgFile)
		singleton = drive.New(api)
	})
	return singleton
}
