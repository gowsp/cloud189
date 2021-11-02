package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			web.GetClient().Ls("/")
			return
		}
		web.GetClient().Ls(args[0])
	},
}
