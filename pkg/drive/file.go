package drive

import (
	"log"
	"path/filepath"

	"github.com/gowsp/cloud189/pkg"
)

type MediaType int

const (
	ALL MediaType = iota
	Pict
	MUSIC
	VIDEO
	DOCUMENT
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
	TB = 1 << 40

	Slice = 10 * MB
)

func (f *Client) parse(path ...string) []pkg.File {
	data := make([]pkg.File, 0)
	for _, path := range path {
		file, err := f.Stat(path)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, file)
	}
	return data
}
func Dir(path string) string {
	dir := filepath.Dir(path)
	return filepath.ToSlash(dir)
}
func IsDir(path string) bool {
	return path[len(path)-1] == '/'
}

func Base(path string) string {
	return filepath.Base(path)
}
