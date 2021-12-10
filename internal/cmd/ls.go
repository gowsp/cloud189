package cmd

import (
	"fmt"
	"os"

	"github.com/gowsp/cloud189-cli/pkg/file"
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			ls("/")
			return
		}
		ls(args[0])
	},
}

func ls(path string) {
	client := web.NewClient(cfgFile)
	info, err := client.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if info.IsDir() {
		files := client.List(info.Id(), 1)
		for _, v := range files {
			fmt.Println(file.ReadableFileInfo(v))
		}
	} else {
		fmt.Println(file.ReadableFileInfo(info))
	}
}
