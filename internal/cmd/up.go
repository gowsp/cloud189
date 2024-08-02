package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var upCfg pkg.UploadConfig

func init() {
	upCmd.Flags().Uint32VarP(&upCfg.Num, "parallel", "p", 5, "number of parallels for file upload")
	upCmd.Flags().StringVarP(&upCfg.Parten, "name", "n", "", "filter filename regular expression")
}

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
		if err := App().Upload(upCfg, cloud, locals...); err != nil {
			fmt.Println(err)
		}
	},
}
