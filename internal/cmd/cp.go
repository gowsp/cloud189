package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		web.GetClient().Cp(args...)
	},
}
