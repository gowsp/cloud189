package file

import (
	"io/fs"

	"github.com/gowsp/cloud189/pkg"
)

func Convert[F pkg.File](loader func() ([]F, error)) (files []fs.FileInfo, err error) {
	data, err := loader()
	if err != nil {
		return nil, err
	}
	for _, file := range data {
		files = append(files, file)
	}
	return
}
