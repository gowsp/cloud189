package cache

import (
	"os"
	"sync"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

type Latest func(id string) pkg.File

type DirEntry struct {
	exp    time.Time
	Info   pkg.File
	dirs   sync.Map
	files  sync.Map
	loaded bool
}

func newEntry(file pkg.File) *DirEntry {
	return &DirEntry{Info: file}
}
func (e *DirEntry) valid() bool {
	return e.loaded && e.exp.After(time.Now())
}
func (e *DirEntry) enable() {
	e.exp = time.Now().Add(time.Minute * 1)
	e.loaded = true
}
func (e *DirEntry) Load(name string) (pkg.File, error) {
	if val, ok := e.dirs.Load(name); ok {
		return val.(*DirEntry).Info, nil
	}
	if val, ok := e.files.Load(name); ok {
		return val.(pkg.File), nil
	}
	return nil, os.ErrNotExist
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
func (e *DirEntry) remove(file pkg.File) {
	if file.IsDir() {
		e.dirs.Delete(file.Name())
	} else {
		e.files.Delete(file.Name())
	}
	e.invalid()
}
func (e *DirEntry) invalid() {
	e.loaded = false
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
