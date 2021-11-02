package cmd

import (
	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login cloud189",
	Run: func(cmd *cobra.Command, args []string) {
		web.GetClient().Login()
	},
}
