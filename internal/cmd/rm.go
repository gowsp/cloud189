package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:    "rm",
	Short:  "remove file",
	PreRun: session.Parse,
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := App().Remove(args...); err != nil {
			fmt.Println(err)
		}
	},
}
