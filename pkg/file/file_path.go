package file

import (
	"log"
	"strings"
)

func CheckPath(paths ...string) {
	for _, v := range paths {
		if !strings.HasPrefix(v, "/") {
			log.Fatalf("path %s must start with /\n", v)
		}
	}
}
