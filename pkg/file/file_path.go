package file

import (
	"container/list"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func CheckPath(paths ...string) {
	for _, v := range paths {
		if !strings.HasPrefix(v, "/") {
			log.Fatalf("path %s must start with /\n", v)
		}
	}
}

type LocalDir struct {
	Folders []string
	Files   map[string]*list.List
	Dirict  *list.List
}

func ReadDir(up string) *LocalDir {
	d := filepath.Dir(up)
	l := len(d)
	files := make(map[string]*list.List)
	filepath.WalkDir(up, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		local := &FilePath{FullPath: path, FileInfo: info}
		dir := filepath.ToSlash(filepath.Dir(path[l+1:]))
		if folder, ok := files[dir]; ok {
			folder.PushBack(local)
			return nil
		}
		f := list.New()
		files[dir] = f
		f.PushBack(local)
		return nil
	})
	var dirict *list.List
	folders := make([]string, 0, len(files))
	for p, v := range files {
		if v.Len() == 0 {
			continue
		}
		if p == "." {
			dirict = v
			continue
		}
		folders = append(folders, p)
	}
	if dirict != nil {
		delete(files, ".")
	}
	return &LocalDir{Folders: folders, Files: files, Dirict: dirict}
}
