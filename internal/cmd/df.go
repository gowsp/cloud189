package cmd

import (
	"fmt"

	"github.com/gowsp/cloud189/pkg/file"
	"github.com/spf13/cobra"
)

var dfCmd = &cobra.Command{
	Use:   "df",
	Short: "show information about the space used",
	Run: func(cmd *cobra.Command, args []string) {
		space, err := App().Space()
		if err != nil {
			fmt.Println(err)
			return
		}
		capacity := space.Capacity
		available := space.Available
		used := capacity - available
		fmt.Printf("%-12s%-12s%-12s%s\n", "Size", "Used", "Avail", "Use%")
		fmt.Printf("%-12s%-12s%-12s%.2f%%\n",
			file.ReadableSize(capacity),
			file.ReadableSize(used),
			file.ReadableSize(available),
			float64(used)*100/float64(capacity),
		)

	},
}
