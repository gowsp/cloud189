package cache

import (
	"os"
	"sync"

	"github.com/gowsp/cloud189/pkg"
)

type Cached func(*DirEntry) bool

type DirEntry struct {
	Info  pkg.File
	Init  bool
	dirs  sync.Map
	files sync.Map
}

func newEntry(file pkg.File) *DirEntry {
	return &DirEntry{Info: file}
}
func (e *DirEntry) Files() []pkg.File {
	data := make([]pkg.File, 0)
	e.dirs.Range(func(key, value interface{}) bool {
		data = append(data, value.(*DirEntry).Info)
		return true
	})
	e.files.Range(func(key, value interface{}) bool {
		data = append(data, value.(pkg.File))
		return true
	})
	return data
}
func (e *DirEntry) Remove(file pkg.File) {
	if file.IsDir() {
		e.dirs.Delete(file.Name())
	} else {
		e.files.Delete(file.Name())
	}
}
func (e *DirEntry) clean() {
	e.dirs = sync.Map{}
	e.files = sync.Map{}
}
func (e *DirEntry) load(name string) (pkg.File, error) {
	if val, ok := e.dirs.Load(name); ok {
		return val.(*DirEntry).Info, nil
	}
	if val, ok := e.files.Load(name); ok {
		return val.(pkg.File), nil
	}
	return nil, os.ErrNotExist
}
