package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:    "mkdir",
	Short:  "mkdir on remote",
	PreRun: session.Parse,
	Args:   cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, arg := range args {
			if err := App().Mkdir(arg); err != nil {
				fmt.Println("mkdir error", arg, err)
			}
		}
	},
}
