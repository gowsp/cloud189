package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var date string
var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Release Date: %s\n", date)
		fmt.Printf("Version: %s\n", version)
	},
}
