package file

import (
	"encoding/json"
	"io/fs"
	"os"
	"sync"

	"github.com/gowsp/cloud189-cli/pkg"
)

type MediaType int

const (
	ALL MediaType = iota
	Pict
	MUSIC
	VIDEO
	DOCUMENT
)

type ReadableFile struct {
	Id       json.Number `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	IsFolder bool        `json:"isFolder,omitempty"`
	Client   pkg.Client
	FileInfo pkg.FileInfo
	loader   sync.Once
	temp     *os.File
}

func (f *ReadableFile) getTemp() *os.File {
	f.loader.Do(func() {
		tempDir := os.TempDir() + "/cloud189/"
		_, err := os.Stat(tempDir)
		if os.IsNotExist(err) {
			os.Mkdir(tempDir, 0755)
		}
		tempFile := tempDir + "/" + f.Id.String() + "_" + f.Name
		f.Client.Get(f.Id.String(), tempFile)
		f.temp, _ = os.OpenFile(tempFile, os.O_CREATE|os.O_RDWR, 0644)
	})
	return f.temp
}

func (f *ReadableFile) Read(p []byte) (n int, err error) {
	if f.FileInfo.IsDir() {
		return 0, os.ErrInvalid
	}
	return f.getTemp().Read(p)
}
func (f *ReadableFile) Write(p []byte) (n int, err error) {
	if f.FileInfo.IsDir() {
		return 0, os.ErrInvalid
	}
	return f.getTemp().Write(p)
}
func (f *ReadableFile) Seek(offset int64, whence int) (int64, error) {
	if f.FileInfo.IsDir() {
		return 0, os.ErrInvalid
	}
	if offset == 0 && whence == 0 && f.temp == nil {
		return 0, nil
	}
	return f.getTemp().Seek(offset, whence)
}
func (f *ReadableFile) Readdir(count int) ([]fs.FileInfo, error) {
	if !f.FileInfo.IsDir() {
		return nil, os.ErrInvalid
	}
	return f.Client.Readdir(f.Id.String(), count), nil
}
func (f *ReadableFile) Stat() (fs.FileInfo, error) {
	if f.FileInfo == nil {
		return f.Client.Stat(f.Id.String())
	}
	return f.FileInfo, nil
}
func (f *ReadableFile) Close() error {
	if f.FileInfo.IsDir() || f.temp == nil {
		return nil
	}
	return f.getTemp().Close()
}
