package term

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/spf13/cobra"
)

var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Print the name of the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(session.Pwd())
	},
}
