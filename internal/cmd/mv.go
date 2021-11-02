package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "move file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		web.GetClient().Mv(args...)
	},
}
