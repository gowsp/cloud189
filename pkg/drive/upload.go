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
func (client *FS) Upload(cfg pkg.UploadConfig, cloud string, locals ...string) error {
	err := cfg.Check()
	if err != nil {
		return err
	}
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
			// u := file.NewFastFile(dir.Id(), local)
			// up = append(up, u)
			continue
		}
		files, err := client.uploadLocal(dir, local)
		if err != nil {
			log.Println(err)
			continue
		}
		up = append(up, files...)
	}
	task := cfg.NewTask()
	uploader := client.api.Uploader()
	for _, v := range up {
		if !cfg.Match(v.Name()) {
			continue
		}
		task.Run(func() {
			if err = uploader.Write(v); err != nil {
				log.Println(err)
			}
		})
	}
	task.Close()
	return nil
}

func (client *FS) uploadLocal(parent pkg.File, local string) ([]pkg.Upload, error) {
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
