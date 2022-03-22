package cmd

import (
	"github.com/gowsp/cloud189/pkg/web"
	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "mkdir on remote",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		web.NewClient(cfgFile).Mkdir(args...)
	},
}
