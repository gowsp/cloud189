package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "mkdir on remote",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := client().Mkdirs(args...)
		if err != nil {
			log.Println(err)
		}
	},
}
