package cmd

import (
	"fmt"
	"log"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		clouds := args[:length-1]
		session.Parse(cmd, clouds)
		err := file.CheckPath(clouds...)
		if err != nil {
			fmt.Println(err)
			return
		}
		local := args[length-1]
		if err := App().Download(local, clouds...); err != nil {
			log.Println(err)
		}
	},
}
