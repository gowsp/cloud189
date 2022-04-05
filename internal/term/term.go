package term

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/session"
	"github.com/peterh/liner"
)

func Start() {
	session.SetWorkDir("/")
	cmd.AddCommand(cdCmd, pwdCmd)

	history := filepath.Join(os.TempDir(), ".cloud189_liner_history")
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(completer)

	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	root := cmd.RootCmd
	for {
		if args, err := line.Prompt(fmt.Sprintf("[cloud189 %s]$ ", session.Base())); err == nil {
			root.SetArgs(strings.Split(args, " "))
			root.Execute()
			line.AppendHistory(args)
		} else if err == liner.ErrPromptAborted {
			fmt.Println("Bye")
			break
		} else {
			log.Print("Error reading line: ", err)
			break
		}
	}
	if f, err := os.Create(history); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
