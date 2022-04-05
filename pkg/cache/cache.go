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

func entry(id string) *DirEntry {
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

func List[F pkg.File](parentId string, loader func() ([]F, error)) ([]pkg.File, error) {
	parent := entry(parentId)
	if parent.valid() {
		return parent.Files(), nil
	}
	parent.clean()
	files, err := loader()
	if err != nil {
		return nil, err
	}
	addAll(parent, files)
	parent.enable()
	return parent.Files(), nil
}
func Find[F pkg.File](parentId, name string, loader func() ([]F, error)) (pkg.File, error) {
	parent := entry(parentId)
	if val, err := parent.load(name); err == nil {
		return val, nil
	}
	files, err := loader()
	if err != nil {
		return nil, err
	}
	addAll(parent, files)
	return parent.load(name)
}
func addAll[F pkg.File](parent *DirEntry, files []F) {
	for _, file := range files {
		addFile(parent, file)
	}
}
func addFile[F pkg.File](parent *DirEntry, file F) {
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
		if node, ok := Cache.nodes.Load(file.PId()); ok {
			if file.IsDir() {
				if node, ok := Cache.nodes.Load(file.Id()); ok {
					delete(node.(*DirEntry))
				}
			}
			node.(*DirEntry).remove(file)
		}
	}
}
func delete(entry *DirEntry) {
	Cache.nodes.Delete(entry.Info.Id())
	entry.dirs.Range(func(key, value interface{}) bool {
		delete(value.(*DirEntry))
		return true
	})
}
func Invalid(files ...pkg.File) {
	for _, file := range files {
		if node, ok := Cache.nodes.Load(file.PId()); ok {
			node.(*DirEntry).remove(file)
		}
	}
}
func InvalidId(id string) {
	if node, ok := Cache.nodes.Load(id); ok {
		node.(*DirEntry).invalid()
	}
}
