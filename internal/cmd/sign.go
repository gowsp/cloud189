package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign",
	Run: func(cmd *cobra.Command, args []string) {
		if err := client().Sign(); err != nil {
			fmt.Println(err)
		}
	},
}
