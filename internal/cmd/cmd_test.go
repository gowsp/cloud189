package cmd

import (
	"testing"
)

func TestLogin(t *testing.T) {
	// rootCmd.SetArgs([]string{"login", "-i", "xxxxx", "xxxxx"})
	rootCmd.SetArgs([]string{"login"})
	rootCmd.Execute()
}
func TestSign(t *testing.T) {
	rootCmd.SetArgs([]string{"sign"})
	rootCmd.Execute()
}
func TestMkdir(t *testing.T) {
	rootCmd.SetArgs([]string{"mkdir", "/demo", "/demo/1/2", "/demo/1/3"})
	rootCmd.Execute()
}
func TestUp(t *testing.T) {
	rootCmd.SetArgs([]string{"up", "https://github.com/gowsp/cloud189/releases/download/v0.2/cloud189_0.2_windows_amd64.zip",
		"../../internal", "../../LICENSE", "/demo/1/2"})
	rootCmd.Execute()
}
func TestLs(t *testing.T) {
	rootCmd.SetArgs([]string{"ls", "/demo/1/2"})
	rootCmd.Execute()
}
func TestDownFile(t *testing.T) {
	rootCmd.SetArgs([]string{"dl", "/demo", "d:/"})
	rootCmd.Execute()
}
func TestDownDir(t *testing.T) {
	rootCmd.SetArgs([]string{"dl", "/demo/", "d:/"})
	rootCmd.Execute()
}
func TestCp(t *testing.T) {
	rootCmd.SetArgs([]string{"cp", "/demo/1/2", "/demo"})
	rootCmd.Execute()
}
func TestMv(t *testing.T) {
	rootCmd.SetArgs([]string{"mv", "/demo/1/3", "/demo"})
	rootCmd.Execute()
}
func TestRm(t *testing.T) {
	rootCmd.SetArgs([]string{"rm", "/demo"})
	rootCmd.Execute()
}
func TestDf(t *testing.T) {
	rootCmd.SetArgs([]string{"df"})
	rootCmd.Execute()
}
func TestWebDav(t *testing.T) {
	rootCmd.SetArgs([]string{"webdav", ":80"})
	rootCmd.Execute()
}
