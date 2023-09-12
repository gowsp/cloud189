package drive

import (
	"io/fs"

	"github.com/gowsp/cloud189/pkg"
)

func (f *FS) NewFile(info pkg.File) fs.File {
	if info.IsDir() {
		return &DirFile{info: info}
	}
	return &File{info: info}
}

type File struct {
	info pkg.File
}

func (a *File) Stat() (fs.FileInfo, error) { return a.info, nil }
func (a *File) Read([]byte) (int, error)   { return 0, nil }
func (a *File) Close() error               { return nil }

type DirFile struct {
	info pkg.File
}

func (a *DirFile) Stat() (fs.FileInfo, error)           { return a.info, nil }
func (a *DirFile) Read([]byte) (int, error)             { return 0, nil }
func (a *DirFile) Close() error                         { return nil }
func (a *DirFile) ReadDir(n int) ([]fs.DirEntry, error) { return nil, nil }
