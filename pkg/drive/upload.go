package drive

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (client *Client) Upload(cloud string, locals ...string) error {
	dir, err := client.Stat(cloud)
	if len(locals) > 1 || os.IsNotExist(err) {
		client.Mkdir(cloud[1:], true)
		dir, _ = client.Stat(cloud)
	}
	for _, local := range locals {
		if IsNetFile(local) {
			f := NewNetFile(dir.Id(), local, client.api)
			f.Upload()
			continue
		}
		if IsFastFile(local) {
			f := NewFastFile(dir.Id(), local, client.api)
			f.Upload()
			continue
		}
		client.uploadLocal(dir.Id(), cloud, local)
	}
	return nil
}

func (client *Client) uploadLocal(parentId, cloud, local string) error {
	dirs := make([]string, 0)
	files := make([]FilePath, 0)
	filepath.WalkDir(local, func(path string, d fs.DirEntry, err error) error {
		rel, err := filepath.Rel(local, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			dirs = append(dirs, rel)
		} else {
			info, err := d.Info()
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(filepath.Dir(rel))
			files = append(files, FilePath{CloudPath: rel, LocalPath: path, FileInfo: info})
		}
		return err
	})
	dir, err := client.api.Mkdirs(parentId, dirs...)
	if err != nil {
		return err
	}
	dir["."] = parentId
	for _, f := range files {
		i := NewLocalFile(dir[f.CloudPath].(string), &f, client.api)
		i.Upload()
	}
	return nil
}
