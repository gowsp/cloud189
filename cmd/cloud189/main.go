package main

import (
	"flag"
	"os"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/term"
)

func main() {
	os.Setenv("EXE_MODE", "1")
	flag.Parse()
	cmd.AddCommand(versionCmd)
	if len(flag.Args()) == 0 {
		term.Start()
	} else {
		cmd.Execute()
	}
}
