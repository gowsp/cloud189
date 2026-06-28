package cmd

import (
	"fmt"
	"io/fs"
	"sort"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var recursiveList bool

func init() {
	lsCmd.Flags().BoolVarP(&recursiveList, "recursive", "r", false, "list recursively")
}

type dirReader interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}

var lsCmd = &cobra.Command{
	Use:    "ls [path]",
	PreRun: session.Parse,
	Short:  "list file",
	Args:   cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := file.CheckPath(args...)
		if err != nil {
			return err
		}
		name := "/"
		if session.Pwd() != "" {
			name = session.Pwd()
		}
		if len(args) == 0 {
		} else if len(args) > 0 {
			name = args[0]
		}
		client := App()
		if jsonOutput {
			entries, err := listJSONEntries(client, name, recursiveList)
			if err != nil {
				return err
			}
			return writeJSON(entries)
		}
		files, err := client.ReadDir(name)
		if err != nil {
			return err
		}
		for _, v := range files {
			info, _ := v.Info()
			fmt.Println(file.ReadableFileInfo(info))
		}
		return nil

	},
}

func listJSONEntries(client dirReader, root string, recursive bool) ([]JSONFileEntry, error) {
	entries, err := listJSONEntriesFrom(client, cleanCloudPath(root), cleanCloudPath(root), recursive)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	return entries, nil
}

func listJSONEntriesFrom(client dirReader, current, base string, recursive bool) ([]JSONFileEntry, error) {
	files, err := client.ReadDir(current)
	if err != nil {
		return nil, err
	}
	entries := make([]JSONFileEntry, 0, len(files))
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			return nil, err
		}
		entryPath := joinCloudPath(current, info.Name())
		if recursive && info.IsDir() {
			childEntries, err := listJSONEntriesFrom(client, entryPath, base, recursive)
			if err != nil {
				return nil, err
			}
			entries = append(entries, childEntries...)
			continue
		}
		entries = append(entries, fileToJSONEntryWithBase(entryPath, base, info))
	}
	return entries, nil
}
