package cmd

import (
	"fmt"
	"os"

	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var confirm bool

var logoutCmd = &cobra.Command{
	Use:          "logout",
	Short:        "logout cloud189",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if confirm {
			return logout()
		}
		liner := liner.NewLiner()
		defer liner.Close()
		reply, err := liner.Prompt("Are you sure to logout? (y/n) ")
		if err != nil {
			return err
		}
		switch reply {
		case "y", "Y":
			return logout()
		}
		return nil
	},
}

func logout() error {
	if cfgFile == "" {
		cfgFile = invoker.DefaultPath()
	}
	err := os.Remove(cfgFile)
	if os.IsNotExist(err) {
		err = nil
	}
	if err == nil {
		fmt.Println("logout success")
	}
	return err
}

func init() {
	logoutCmd.Flags().BoolVarP(&confirm, "f", "f", false, "log out now")
}
