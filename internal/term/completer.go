package term

import (
	"path"
	"strings"

	"github.com/gowsp/cloud189/internal/cmd"
	"github.com/gowsp/cloud189/internal/session"
)

var tips = map[string]struct{}{
	"cd":    {},
	"ls":    {},
	"mkdir": {},
	"cp":    {},
	"mv":    {},
	"rm":    {},
	"dl":    {},
	"up":    {},
}
var cmds = []string{"cd", "ls", "mkdir", "cp", "mv", "rm", "dl", "up",
	"pwd", "version", "login", "exit", "logout"}

func completer(line string) (c []string) {
	args := strings.Split(line, " ")
	len := len(args)
	if len < 2 {
		for _, n := range cmds {
			if strings.HasPrefix(n, line) {
				c = append(c, n)
			}
		}
		return
	}
	if _, ok := tips[args[0]]; !ok {
		return
	}
	arg := args[len-1]
	dir, n := path.Split(arg)
	name := session.Join(dir)
	files, err := cmd.App().ListBy(name)
	if err != nil {
		return
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), n) {
			args[len-1] = path.Join(dir, file.Name())
			val := strings.Join(args, " ")
			c = append(c, val)
		}
	}
	return
}
