package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189-cli/pkg/web"
	"github.com/spf13/cobra"
)

var UsePwd bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login cloud189",
	Args: func(cmd *cobra.Command, args []string) error {
		if UsePwd && len(args) < 2 {
			return fmt.Errorf("requires username password parameter, received %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if UsePwd {
			web.NewContent().PwdLogin(args[0], args[1])
			return
		}
		web.NewContent().QrLogin()
	},
}

func init() {
	loginCmd.Flags().BoolVarP(&UsePwd, "i", "i", false, "use username and password to login")
}
