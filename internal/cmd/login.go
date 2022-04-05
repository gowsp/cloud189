package cmd

import (
	"fmt"

	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

var usePwd bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login cloud189",
	Args: func(cmd *cobra.Command, args []string) error {
		if usePwd && len(args) < 2 {
			return fmt.Errorf("requires username password parameter, received %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if usePwd {
			loginFunc(args[0], args[1])
			return
		}
		line := liner.NewLiner()
		defer line.Close()
		username, _ := line.Prompt("username: ")
		password, _ := line.PasswordPrompt("password: ")
		loginFunc(username, password)
	},
}

func loginFunc(username, password string) {
	if err := App().Login(username, password); err != nil {
		fmt.Printf("\n%s\n", err)
		return
	}
	fmt.Println("login success")
}

func init() {
	loginCmd.Flags().BoolVarP(&usePwd, "i", "i", false, "input username and password to login")
}
