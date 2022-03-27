package cmd

import (
	"github.com/spf13/cobra"
)

var webdavCmd = &cobra.Command{
	Use:   "webdav",
	Short: "start webdav server, arg: port",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// webdav.Serve(args[0], web.NewClient(cfgFile))
	},
}
