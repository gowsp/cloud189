package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		web.GetClient().Rm(args...)
	},
}
