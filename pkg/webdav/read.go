package webdav

import (
	"io/fs"
	"os"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"golang.org/x/net/webdav"
)

func newRead(app pkg.App, name string) (webdav.File, error) {
	stat, err := app.Stat(name)
	if err != nil {
		return nil, err
	}
	return &read{app: app, stat: stat}, nil
}

type read struct {
	app  pkg.App
	stat pkg.File
	load sync.Once
	temp *os.File
}

func (r *read) getTemp() *os.File {
	r.load.Do(func() {
		tempDir := os.TempDir() + "/cloud189/"
		_, err := os.Stat(tempDir)
		if os.IsNotExist(err) {
			os.Mkdir(tempDir, 0755)
		}
		tempFile := tempDir + "/" + r.stat.Id() + "_" + r.stat.Name()
		r.app.DownloadFile(tempFile, r.stat)
		r.temp, _ = os.OpenFile(tempFile, os.O_CREATE|os.O_RDWR, 0644)
	})
	return r.temp
}
func (r *read) Seek(offset int64, whence int) (int64, error) { return r.getTemp().Seek(offset, whence) }
func (r *read) Read(p []byte) (n int, err error)             { return r.getTemp().Read(p) }
func (r *read) Write(p []byte) (n int, err error)            { return r.getTemp().Write(p) }
func (r *read) Close() error                                 { return nil }
func (r *read) Readdir(count int) ([]fs.FileInfo, error) {
	if !r.stat.IsDir() {
		return nil, os.ErrInvalid
	}
	file, err := r.app.List(r.stat)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0)
	for _, v := range file {
		infos = append(infos, v)
	}
	return infos, nil
}
func (r *read) Stat() (info fs.FileInfo, err error) {
	return r.stat, err
}
