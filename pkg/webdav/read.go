package webdav

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"golang.org/x/net/webdav"
)

func newRead(app pkg.Drive, name string) (webdav.File, error) {
	stat, err := app.Stat(name)
	if err != nil {
		return nil, err
	}
	return &read{app: app, name: name, stat: stat.(pkg.File)}, nil
}

var empty = &read{}

type read struct {
	app  pkg.Drive
	name string
	stat pkg.File
	load sync.Once
	temp *os.File
}

func (r *read) getTemp() *os.File {
	r.load.Do(func() {
		dir, name := filepath.Split(r.name)
		tempDir := os.TempDir() + "/cloud189" + dir
		_, err := os.Stat(tempDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(tempDir, 0755)
			if err != nil {
				log.Println(err)
				return
			}
		}
		r.app.Download(tempDir, r.name)
		r.temp, _ = os.OpenFile(tempDir+"/"+name, os.O_CREATE|os.O_RDWR, 0644)
	})
	return r.temp
}
func (r *read) Seek(offset int64, whence int) (int64, error) { return r.getTemp().Seek(offset, whence) }
func (r *read) Read(p []byte) (n int, err error)             { return r.getTemp().Read(p) }
func (r *read) Write(p []byte) (n int, err error)            { return r.getTemp().Write(p) }
func (r *read) Close() error                                 { return nil }
func (r *read) Readdir(count int) ([]fs.FileInfo, error) {
	data, err := r.app.ReadDir(r.name)
	if err != nil {
		return nil, err
	}
	files := make([]fs.FileInfo, len(data))
	for i, v := range data {
		files[i], _ = v.Info()
	}
	return files, err
}
func (r *read) Stat() (info fs.FileInfo, err error) {
	return r.stat, err
}
