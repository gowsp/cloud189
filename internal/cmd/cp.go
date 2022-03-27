package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		dest := args[length-1]
		from := args[:length-1]
		err := client().Copy(dest, from...)
		if err != nil {
			fmt.Println(err)
		}
	},
}
