package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		clouds := args[:length-1]
		local := args[length-1]
		err := client().Download(local, clouds...)
		if err != nil {
			log.Println(err)
		}
	},
}
