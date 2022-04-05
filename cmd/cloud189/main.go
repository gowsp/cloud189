package main

import (
	"flag"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/term"
)

func main() {
	flag.Parse()
	cmd.AddCommand(versionCmd)
	if len(flag.Args()) == 0 {
		term.Start()
	} else {
		cmd.Execute()
	}
}
