package cmd

import (
	"github.com/gowsp/cloud189/pkg/web"
	"github.com/spf13/cobra"
)

var dfCmd = &cobra.Command{
	Use:   "df",
	Short: "show information about the space used",
	Run: func(cmd *cobra.Command, args []string) {
		web.NewClient(cfgFile).Df()
	},
}
