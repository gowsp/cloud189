package cmd

import (
	"fmt"
	"os"

	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := client()
		info, err := client.Stat(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if info.IsDir() {
			files, err := client.List(info)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, v := range files {
				fmt.Println(file.ReadableFileInfo(v))
			}
			return
		}
		fmt.Println(file.ReadableFileInfo(info))
	},
}
