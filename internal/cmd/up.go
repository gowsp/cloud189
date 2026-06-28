package cmd

import (
	"errors"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var upCfg pkg.UploadConfig
var upInput string
var upRemotePath string

func init() {
	upCmd.Flags().Uint32VarP(&upCfg.Num, "parallel", "p", 5, "number of parallels for file upload")
	upCmd.Flags().StringVarP(&upCfg.Parten, "name", "n", "", "filter filename regular expression")
	upCmd.Flags().StringVar(&upInput, "input", "", "local input file path for machine-mode upload")
	upCmd.Flags().StringVar(&upRemotePath, "path", "", "full remote file path for machine-mode upload")
	upCmd.Flags().StringVar(&upCfg.Policy, "policy", "skip", "upload policy for machine-mode upload: skip or overwrite")
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonOutput || upInput != "" || upRemotePath != "" {
			if len(args) != 0 {
				return cobra.NoArgs(cmd, args)
			}
			if upInput == "" {
				return errors.New("--input is required in machine mode")
			}
			if upRemotePath == "" {
				return errors.New("--path is required in machine mode")
			}
			paths := []string{upRemotePath}
			session.Parse(cmd, paths)
			if err := file.CheckPath(paths[0]); err != nil {
				return err
			}
			info, err := App().UploadFile(upCfg, upInput, paths[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeJSON(fileToJSONEntry(paths[0], info))
			}
			return nil
		}
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return err
		}
		length := len(args)
		cloud := session.Join(args[length-1])
		err := file.CheckPath(cloud)
		if err != nil {
			return err
		}
		locals := args[:length-1]
		if err := App().Upload(upCfg, cloud, locals...); err != nil {
			return err
		}
		return nil
	},
}
