package drive

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

func (client *FS) UploadFrom(file pkg.Upload) error {
	uploader := client.api.Uploader()
	return uploader.Write(file)
}
func (client *FS) Upload(cloud string, locals ...string) error {
	dir, err := client.stat(cloud)
	if len(locals) > 1 || os.IsNotExist(err) {
		client.Mkdir(cloud[1:])
		dir, _ = client.stat(cloud)
	}
	up := make([]pkg.Upload, 0)
	for _, local := range locals {
		if file.IsNetFile(local) {
			up = append(up, file.NewURLFile(dir.Id(), local))
			continue
		}
		if file.IsFastFile(local) {
			u := file.NewFastFile(dir.Id(), local)
			up = append(up, u)
			continue
		}
		files, _ := client.uploadLocal(dir, cloud, local)
		up = append(up, files...)
	}
	uploader := client.api.Uploader()
	for _, v := range up {
		err := uploader.Write(v)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (client *FS) uploadLocal(parent pkg.File, cloud, local string) ([]pkg.Upload, error) {
	stat, err := os.Stat(local)
	if err != nil {
		return nil, err
	}
	up := make([]pkg.Upload, 0)
	if !stat.IsDir() {
		up = append(up, file.NewLocalFile(parent.Id(), local))
		return up, nil
	}
	dirs := map[string]string{
		".": parent.Id(),
	}
	filepath.WalkDir(local, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			rel := file.Rel(local, path)
			if rel == "." {
				return nil
			}
			f, _ := client.api.Mkdir(parent, rel)
			dirs[rel] = f.Id()
			return nil
		}
		dir, _ := filepath.Split(path)
		rel := file.Rel(local, dir)
		up = append(up, file.NewLocalFile(dirs[rel], path))
		return err
	})
	return up, nil
}
