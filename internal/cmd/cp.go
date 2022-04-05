package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:    "cp",
	Short:  "copy file",
	PreRun: session.Parse,
	Args:   cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		length := len(args)
		dest := args[length-1]
		from := args[:length-1]
		if err = App().Copy(dest, from...); err != nil {
			fmt.Println(err)
		}
	},
}
