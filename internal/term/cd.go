package term

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/session"
	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:    "cd",
	Short:  "change dir",
	Args:   cobra.MaximumNArgs(1),
	PreRun: session.Parse,
	Run: func(command *cobra.Command, args []string) {
		if len(args) == 0 {
			session.SetWorkDir("/")
			return
		}
		path := args[0]
		stat, err := cmd.App().Stat(path)
		if err == nil && stat.IsDir() {
			session.SetWorkDir(path)
			return
		}
		fmt.Printf("cd: %s: Not a directory\n", path)
	},
}
