package drive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

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
