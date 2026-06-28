package main

import (
	"os"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/term"
)

func main() {
	os.Setenv("EXE_MODE", "1")
	cmd.AddCommand(versionCmd)
	if len(os.Args) == 1 {
		term.Start()
	} else {
		cmd.Execute()
	}
}
