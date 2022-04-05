package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		length := len(args)
		cloud := session.Join(args[length-1])
		err := file.CheckPath(cloud)
		if err != nil {
			fmt.Println(err)
			return
		}
		locals := args[:length-1]
		App().Upload(cloud, locals...)
	},
}
