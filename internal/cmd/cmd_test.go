package cmd

import (
	"testing"
)

func TestLs(t *testing.T) {
	rootCmd.SetArgs([]string{"ls", "/demo"})
	rootCmd.Execute()
}
func TestRm(t *testing.T) {
	rootCmd.SetArgs([]string{"rm", "/demo"})
	rootCmd.Execute()
}