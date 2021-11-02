package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "cloud189",
		Long: "cloud189 enables users to manage cloud files through the command line. For more information, please visit https://github.com/gowsp/cloud189",
	}
)

func Execute(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(dlCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(mkdirCmd)
	rootCmd.AddCommand(mvCmd)
	rootCmd.AddCommand(cpCmd)
}
