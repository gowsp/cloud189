package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/pkg/app"
	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			cfgFile = invoker.DefaultPath()
		}
		app := app.New(cfgFile)
		if err := app.Sign(); err != nil {
			fmt.Println(err)
		}
	},
}
