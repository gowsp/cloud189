package cmd

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:    "stat [path]",
	Short:  "show file metadata",
	PreRun: session.Parse,
	Args:   cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "/"
		if session.Pwd() != "" {
			name = session.Pwd()
		}
		if len(args) > 0 {
			name = args[0]
		}
		if err := file.CheckPath(name); err != nil {
			return err
		}
		info, err := App().Stat(name)
		if err != nil {
			if jsonOutput && errors.Is(err, fs.ErrNotExist) {
				return writeJSON(nil)
			}
			return err
		}
		if jsonOutput {
			return writeJSON(fileToJSONEntry(name, info))
		}
		fmt.Println(file.ReadableFileInfo(info))
		return nil
	},
}
