package cmd

import (
	"errors"
	"log"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var dlOutput string

func init() {
	dlCmd.Flags().StringVarP(&dlOutput, "output", "o", "", "local output file path for machine-mode download")
}

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "download file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonOutput || dlOutput != "" {
			if len(args) != 1 {
				return cobra.ExactArgs(1)(cmd, args)
			}
			session.Parse(cmd, args)
			if err := file.CheckPath(args[0]); err != nil {
				return err
			}
			if dlOutput == "" {
				return errors.New("--output is required in machine mode")
			}
			info, err := App().DownloadTo(dlOutput, args[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeJSON(fileToJSONEntry(args[0], info))
			}
			return nil
		}
		length := len(args)
		if length < 2 {
			return cobra.MinimumNArgs(2)(cmd, args)
		}
		clouds := args[:length-1]
		session.Parse(cmd, clouds)
		err := file.CheckPath(clouds...)
		if err != nil {
			return err
		}
		local := args[length-1]
		if err := App().Download(local, clouds...); err != nil {
			log.Println(err)
		}
		return nil
	},
}
