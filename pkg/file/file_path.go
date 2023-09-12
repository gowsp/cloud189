package file

import (
	"fmt"
	"path"
	"path/filepath"
)

func Rel(parent, file string) string {
	rel, _ := filepath.Rel(parent, file)
	return filepath.ToSlash(rel)
}

func CheckPath(paths ...string) error {
	for _, v := range paths {
		if !path.IsAbs(v) {
			return fmt.Errorf("path %s must start with /", v)
		}
	}
	return nil
}
