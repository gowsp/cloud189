package file

import (
	"fmt"
	"path"
)

func CheckPath(paths ...string) error {
	for _, v := range paths {
		if !path.IsAbs(v) {
			return fmt.Errorf("path %s must start with /", v)
		}
	}
	return nil
}
