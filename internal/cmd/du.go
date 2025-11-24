package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/internal/session"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var duCmd = &cobra.Command{
	Use:    "du",
	Short:  "show file usage statistics",
	PreRun: session.Parse,
	Args:   cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := file.CheckPath(args...)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		var path string
		if len(args) == 0 {
			path = session.Pwd()
		} else {
			path = args[0]
		}
		
		// First check if path is a directory
		fileInfo, err := App().Stat(path)
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}
		
		if fileInfo.IsDir() {
			// If it's a directory, list the files inside and show usage for each
			files, err := App().ReadDir(path)
			if err != nil {
				fmt.Println("Error listing directory:", err)
				return
			}
			
			// Print header
			fmt.Printf("%-10s %-10s %-10s %s\n", "FILES", "FOLDERS", "SIZE", "NAME")
			
			totalFiles := uint64(0)
			totalSize := uint64(0)
			totalFolders := uint64(0)
			
			for _, v := range files {
				info, _ := v.Info()
				usage, err := App().Usage(path + "/" + info.Name())
				if err != nil {
					fmt.Printf("Error getting usage for %s: %v\n", info.Name(), err)
					continue
				}
				
				// Display usage in a list format
				fmt.Printf("%-10d %-10d %-10s %s\n", 
					usage.FileCount(), 
					usage.FolderCount(), 
					file.ReadableSize(usage.FileSize()),
					info.Name())
				
				totalFiles += usage.FileCount()
				totalSize += usage.FileSize()
				totalFolders += usage.FolderCount()
			}
			
			// Print total
			fmt.Printf("%-10d %-10d %-10s %s\n", totalFiles, totalFolders, file.ReadableSize(totalSize), "Total")
		} else {
			// If it's a file, show usage statistics
			usage, err := App().Usage(path)
			if err != nil {
				fmt.Println("Error getting usage:", err)
				return
			}
			
			fmt.Printf("%-10s %-10s %-10s %s\n", "FILES", "FOLDERS", "SIZE", "NAME")
			fmt.Printf("%-10d %-10d %-10s %s\n", 
				usage.FileCount(), 
				usage.FolderCount(), 
				file.ReadableSize(usage.FileSize()),
				fileInfo.Name())
		}
	},
}