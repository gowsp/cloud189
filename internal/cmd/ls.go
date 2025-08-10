package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var jsonFormat bool

func init() {
	lsCmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "print json format")
}

type FileInfoEntry struct {
	Name  string
	IsDir bool
	Type  fs.FileMode
	Info  os.FileInfo
}

var lsCmd = &cobra.Command{
	Use:    "ls",
	PreRun: session.Parse,
	Short:  "list file",
	Args:   cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		var name string
		if len(args) == 0 {
			name = session.Pwd()
		} else {
			name = args[0]
		}
		client := App()
		files, err := client.ReadDir(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		if jsonFormat {
			entries := make([]FileInfoEntry, len(files))
			for i, f := range files {
				info, _ := f.Info()
				entries[i] = FileInfoEntry{
					Name:  info.Name(),
					IsDir: info.IsDir(),
					Type:  info.Mode().Type(),
					Info:  info,
				}
			}
			data, err := json.Marshal(entries)
			if err != nil {
				fmt.Print("{\"error\":\"", err.Error(), "\"}")
			} else {
				fmt.Print(string(data))
			}
		} else {
			for _, v := range files {
				info, _ := v.Info()
				fmt.Println(file.ReadableFileInfo(info))
			}
		}

	},
}
