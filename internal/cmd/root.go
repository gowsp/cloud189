package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/app"
	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	jsonOutput bool
	RootCmd    = &cobra.Command{
		Use:           "cloud189",
		Long:          "cloud189 是一个基于天翼云接口的命令行客户端。详情请访问 https://github.com/gowsp/cloud189",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func AddCommand(cmds ...*cobra.Command) {
	RootCmd.AddCommand(cmds...)
}

func ResetTermFlags() {
	jsonOutput = false
	recursiveList = false
	dlOutput = ""
	upInput = ""
	upRemotePath = ""
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if jsonOutput {
			_ = writeJSONError(err)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "指定配置文件路径（默认为 $HOME/.config/cloud189/config.json）")
	RootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "以 JSON 格式输出结果")

	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(qrLoginCmd)
	RootCmd.AddCommand(logoutCmd)
	RootCmd.AddCommand(signCmd)
	RootCmd.AddCommand(upCmd)
	RootCmd.AddCommand(rmCmd)
	RootCmd.AddCommand(dlCmd)
	RootCmd.AddCommand(lsCmd)
	RootCmd.AddCommand(statCmd)
	RootCmd.AddCommand(mkdirCmd)
	RootCmd.AddCommand(mvCmd)
	RootCmd.AddCommand(cpCmd)
	RootCmd.AddCommand(dfCmd)
	RootCmd.AddCommand(duCmd)
	RootCmd.AddCommand(webdavCmd)
	RootCmd.AddCommand(shareCmd)
	RootCmd.AddCommand(saveCmd)
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
