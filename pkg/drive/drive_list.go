package drive

import (
	"io/fs"

	"github.com/gowsp/cloud189/pkg"
)

func (f *FS) Stat(name string) (fs.FileInfo, error) {
	return f.stat(name)
}
func (f *FS) resolve(path ...string) (files []pkg.File) {
	for _, path := range path {
		file, err := f.stat(path)
		if err != nil {
			continue
		}
		files = append(files, file)
	}
	return
}
func (f *FS) stat(name string) (pkg.File, error) {
	if name == "/" {
		return f.root, nil
	}
	return f.api.Stat(name)
	// var err error
	// var file pkg.File = f.root
	// path := strings.Split(name, "/")
	// size := len(path) - 1
	// for i := 1; i < size; i++ {
	// 	file, err = f.search(file, pkg.DIR, path[i])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// if path[size] == "" {
	// 	return file, nil
	// }
	// return f.search(file, pkg.ALL, path[size])
}
func (f *FS) search(parent pkg.File, fileType pkg.FileType, name string) (pkg.File, error) {
	return load(parent.Id()).search(name, func() (pkg.File, error) {
		entry, err := f.api.Search(parent, fileType, name)
		if err != nil {
			return nil, err
		}
		for _, file := range entry {
			if file.Name() == name {
				return file, nil
			}
		}
		return nil, fs.ErrNotExist
	})
}

func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	dir, err := f.stat(name)
	if err != nil {
		return nil, err
	}
	node := load(dir.Id())
	if node == nil {
		node = newNode(dir)
	}
	info, err := node.list(func() ([]pkg.File, error) {
		return f.api.List(dir, pkg.ALL)
	})
	if err != nil {
		return nil, err
	}
	result := make([]fs.DirEntry, 0, len(info))
	for _, v := range info {
		result = append(result, v.(fs.DirEntry))
	}
	return result, nil
}
