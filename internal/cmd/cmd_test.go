package cmd

import (
	"testing"
)

func TestLogin(t *testing.T) {
	// rootCmd.SetArgs([]string{"login", "-i", "xxxxx", "xxxxx"})
	RootCmd.SetArgs([]string{"login"})
	RootCmd.Execute()
}
func TestSign(t *testing.T) {
	RootCmd.SetArgs([]string{"sign"})
	RootCmd.Execute()
}
func TestMkdir(t *testing.T) {
	RootCmd.SetArgs([]string{"mkdir", "/demo", "/demo/1", "/demo/1/2", "/demo/1/3"})
	RootCmd.Execute()
}
func TestUp(t *testing.T) {
	RootCmd.SetArgs([]string{"up",
		"cmd_test.go",
		"fast://B4CC2601D293ECB814447B80C0ACEC8D:3071418/cloud189_fast.zip",
		"https://github.com/gowsp/cloud189/releases/download/v0.4.3/cloud189_0.4.3_windows_amd64.zip",
		"../../internal",
		 "/demo/1",
	})
	RootCmd.Execute()
}
func TestLs(t *testing.T) {
	RootCmd.SetArgs([]string{"ls", "/demo/1/2"})
	RootCmd.Execute()
}
func TestDownFile(t *testing.T) {
	RootCmd.SetArgs([]string{"dl", "/demo/1/2/LICENSE", "/tmp/"})
	RootCmd.Execute()
}
func TestDownDir(t *testing.T) {
	RootCmd.SetArgs([]string{"dl", "/demo/", "/tmp/"})
	RootCmd.Execute()
}
func TestCp(t *testing.T) {
	RootCmd.SetArgs([]string{"cp", "/demo/1/2", "/demo"})
	RootCmd.Execute()
}
func TestMv(t *testing.T) {
	RootCmd.SetArgs([]string{"mv", "/demo/1/3", "/demo"})
	RootCmd.Execute()
}
func TestRm(t *testing.T) {
	RootCmd.SetArgs([]string{"rm", "/demo"})
	RootCmd.Execute()
}
func TestDf(t *testing.T) {
	RootCmd.SetArgs([]string{"df"})
	RootCmd.Execute()
}
func TestWebDav(t *testing.T) {
	RootCmd.SetArgs([]string{"webdav", ":8080"})
	RootCmd.Execute()
}
