package cmd

import (
	"github.com/gowsp/cloud189/pkg/web"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign",
	Run: func(cmd *cobra.Command, args []string) {
		web.NewClient(cfgFile).Sign()
	},
}
