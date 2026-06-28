package drive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/gowsp/cloud189/pkg"
)

func (f *FS) Download(local string, cloud ...string) error {
	info, err := os.Stat(local)
	if err != nil {
		return err
	}
	sources := f.resolve(cloud...)
	if len(sources) > 0 && !info.IsDir() {
		return errors.New("local param need dir")
	}
	for _, source := range sources {
		if err = f.download(info, local, source); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (f *FS) DownloadTo(localPath, remoteFilePath string) (os.FileInfo, error) {
	source, err := f.stat(remoteFilePath)
	if err != nil {
		return nil, err
	}
	if source.IsDir() {
		return nil, errors.New("not support download dir")
	}
	if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return nil, err
	}
	d, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	resp, err := f.api.Download(source, 0)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, errors.New("error download status code " + resp.Status)
	}
	defer resp.Body.Close()
	if _, err = io.Copy(d, resp.Body); err != nil {
		return nil, err
	}
	info, err := d.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() != source.Size() {
		return nil, fmt.Errorf("download size mismatch: got %d, want %d", info.Size(), source.Size())
	}
	return source, nil
}

func (f *FS) download(info os.FileInfo, local string, source pkg.File) error {
	if info.IsDir() {
		local = path.Join(local, source.Name())
	}
	d, err := os.OpenFile(local, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer d.Close()
	info, err = d.Stat()
	if info.Size() == source.Size() {
		return nil
	}
	if err != nil {
		return err
	}
	resp, err := f.api.Download(source, info.Size())
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New("error download status code " + resp.Status)
	}
	defer resp.Body.Close()
	io.Copy(d, resp.Body)
	return nil
}
