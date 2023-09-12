package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:    "ls",
	PreRun: session.Parse,
	Short:  "list file",
	Args:   cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		var name string
		if len(args) == 0 {
			name = session.Pwd()
		} else {
			name = args[0]
		}
		client := App()
		files, err := client.ReadDir(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range files {
			info, _ := v.Info()
			fmt.Println(file.ReadableFileInfo(info))
		}
	},
}
