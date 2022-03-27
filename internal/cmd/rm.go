package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := client().Remove(args...); err != nil {
			fmt.Println(err)
		}
	},
}
