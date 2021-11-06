package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		clouds := args[:length-1]
		local := args[length-1]
		web.GetClient().Dl(local, clouds...)
	},
}
