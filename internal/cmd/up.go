package cmd

import (
	"github.com/gowsp/cloud189/pkg/web"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		cloud := args[length-1]
		locals := args[:length-1]
		web.NewClient(cfgFile).Up(cloud, locals...)
	},
}
