package session

import (
	"path"

	"github.com/spf13/cobra"
)

var workingDirectory = ""

func Pwd() string {
	return workingDirectory
}
func Base() string {
	return path.Base(workingDirectory)
}
func SetWorkDir(path string) {
	workingDirectory = path
}
func Join(name string) string {
	if path.IsAbs(name) {
		return name
	}
	return path.Join(workingDirectory, name)
}
func Parse(cmd *cobra.Command, args []string) {
	for i, arg := range args {
		args[i] = Join(arg)
	}
}
