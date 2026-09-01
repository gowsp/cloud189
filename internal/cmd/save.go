package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "save share files to cloud",
	Long: `Save files from a Tianyi Cloud share link to your cloud drive.

Examples:
  cloud189 save https://cloud.189.cn/t/reiQruJJnyEv /target/path
  cloud189 save "https://cloud.189.cn/t/reiQruJJnyEv（访问码：v5oy）" /target/path`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		shareURL := args[0]
		target := args[1]
		if err := App().SaveShare(shareURL, target); err != nil {
			fmt.Println(err)
		}
	},
}
