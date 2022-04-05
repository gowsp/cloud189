package term

import (
	"fmt"
	"testing"

	"github.com/gowsp/cloud189/internal/session"
)

func TestCompleter(t *testing.T) {
	session.SetWorkDir("/")
	ls := completer("cd 我")
	fmt.Println(ls)
}
