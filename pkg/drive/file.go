package drive

import (
	"github.com/gowsp/cloud189/pkg"
)

func (f *Client) parse(path ...string) (files []pkg.File) {
	for _, path := range path {
		file, err := f.Stat(path)
		if err != nil {
			continue
		}
		files = append(files, file)
	}
	return
}
