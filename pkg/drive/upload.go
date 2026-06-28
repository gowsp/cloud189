package drive

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
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
		files, err := client.uploadLocal(dir, local, cfg.Parten)
		if err != nil {
			log.Println(err)
			continue
		}
		up = append(up, files...)
	}
	task := cfg.NewTask()
	uploader := client.api.Uploader()
	for _, v := range up {
		r := v
		task.Run(func() {
			if err = uploader.Write(r); err != nil {
				log.Println(err)
			}
		})
	}
	task.Close()
	return nil
}

func (client *FS) UploadFile(cfg pkg.UploadConfig, localPath, remoteFilePath string) (fs.FileInfo, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	if cfg.Policy == "" {
		cfg.Policy = "skip"
	}
	localInfo, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}
	if localInfo.IsDir() {
		return nil, errors.New("input must be a file")
	}
	remoteFilePath = path.Clean(remoteFilePath)
	if !path.IsAbs(remoteFilePath) {
		return nil, errors.New("remote path must start with /")
	}
	existing, err := client.stat(remoteFilePath)
	if err == nil {
		if existing.IsDir() {
			return nil, errors.New("remote path is a directory")
		}
		if cfg.Policy == "skip" {
			if existing.Size() == localInfo.Size() {
				return existing, nil
			}
			return nil, fmt.Errorf("remote file exists with different size: got %d, want %d", existing.Size(), localInfo.Size())
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	parentPath, remoteName := path.Split(remoteFilePath)
	if remoteName == "" {
		return nil, errors.New("remote path must include a file name")
	}
	parent, err := client.ensureDir(parentPath)
	if err != nil {
		return nil, err
	}
	upload, err := file.NewLocalFileWithOptions(parent.Id(), localPath, remoteName, cfg.Policy == "overwrite")
	if err != nil {
		return nil, err
	}
	if closer, ok := upload.(interface{ Close() }); ok {
		defer closer.Close()
	}
	if err = client.api.Uploader().Write(upload); err != nil {
		return nil, err
	}
	if node := load(parent.Id()); node != nil {
		node.invalid()
	}
	return client.stat(remoteFilePath)
}

func (client *FS) ensureDir(name string) (pkg.File, error) {
	name = path.Clean(name)
	if name == "." || name == "/" {
		return client.root, nil
	}
	dir, err := client.stat(name)
	if err == nil {
		if !dir.IsDir() {
			return nil, errors.New("parent path is not a directory")
		}
		return dir, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	if _, err = client.api.Mkdir(file.Root, name[1:]); err != nil {
		return nil, err
	}
	return client.stat(name)
}

func (client *FS) uploadLocal(parent pkg.File, local string, parten string) ([]pkg.Upload, error) {
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
	parten = filepath.Join(local, parten)
	files, err := filepath.Glob(parten)
	if err != nil {
		return nil, err
	}
	for _, localFile := range files {
		info, err := os.Stat(localFile)
		if err != nil {
			log.Println(err)
			continue
		}
		if info.IsDir() {
			filepath.WalkDir(localFile, func(path string, d fs.DirEntry, err error) error {
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
		} else {
			up = append(up, file.NewLocalFile(parent.Id(), localFile))
		}

	}
	return up, nil
}
