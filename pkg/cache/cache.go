package cache

import (
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

var Cache = &cache{}

type cache struct {
	nodes sync.Map
}

func Entry(id string) *DirEntry {
	if val, ok := Cache.nodes.Load(id); ok {
		return val.(*DirEntry)
	}
	if id == file.Root.Id() {
		entry := newEntry(file.Root)
		Cache.nodes.Store(id, entry)
		return entry
	}
	return nil
}

func List[F pkg.File](parentId string, cached Cached, loader func() ([]F, error)) ([]pkg.File, error) {
	parent := Entry(parentId)
	if cached(parent) {
		return parent.Files(), nil
	}
	files, err := loader()
	if err != nil {
		return nil, err
	}
	parent.clean()
	addAll(parent, files)
	return parent.Files(), nil
}
func Load(parentId, name string, loader func() error) (pkg.File, error) {
	parent := Entry(parentId)
	if val, err := parent.load(name); err == nil {
		return val, nil
	}
	err := loader()
	if err != nil {
		return nil, err
	}
	return parent.load(name)
}
func addAll[F pkg.File](parent *DirEntry, files []F) {
	for _, file := range files {
		AddFile(parent, file)
	}
}
func AddFile[F pkg.File](parent *DirEntry, file F) {
	if !file.IsDir() {
		parent.files.Store(file.Name(), file)
		return
	}
	if dir, ok := parent.dirs.Load(file.Name()); ok {
		dir.(*DirEntry).Info = file
	} else {
		entry := newEntry(file)
		Cache.nodes.Store(file.Id(), entry)
		parent.dirs.Store(file.Name(), entry)
	}
}

func Delete(files ...pkg.File) {
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		if node, ok := Cache.nodes.Load(file.Id()); ok {
			remove(node.(*DirEntry))
		}
	}
}
func remove(entry *DirEntry) {
	Cache.nodes.Delete(entry.Info.Id())
	entry.dirs.Range(func(key, value interface{}) bool {
		remove(value.(*DirEntry))
		return true
	})
}
