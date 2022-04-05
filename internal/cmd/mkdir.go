package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var parents bool
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
		if parents {
			if err := App().Mkdirs(args...); err != nil {
				fmt.Println(err)
			}
			return
		}
		for _, arg := range args {
			if err := App().Mkdir(arg, false); err != nil {
				fmt.Println("mkdir error", arg, err)
			}
		}
	},
}

func init() {
	mkdirCmd.Flags().BoolVarP(&parents, "p", "p", false, "no error if existing, make parent directories as needed")
}
