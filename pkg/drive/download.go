package drive

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gowsp/cloud189/pkg"
)

func (f *Client) Download(local string, paths ...string) error {
	for _, path := range paths {
		file, err := f.Stat(path)
		if err != nil {
			fmt.Printf("%s %s\n", path, err)
			continue
		}
		if file.IsDir() {
			f.downloadDir(local, file)
		} else {
			f.DownloadFile(local, file)
		}
	}
	return nil
}

func (c *Client) downloadDir(local string, cloud pkg.File) error {
	files, err := c.api.ListFile(cloud.Id())
	if err != nil {
		return err
	}
	local = path.Clean(local + "/" + cloud.Name())
	if _, err = os.Stat(local); os.IsNotExist(err) {
		if err = os.Mkdir(local, 0766); err != nil {
			fmt.Println(err)
			return err
		}
	}
	for _, f := range files {
		if f.IsDir() {
			c.downloadDir(local, f)
		} else {
			c.DownloadFile(local, f)
		}
	}
	return nil
}
func (c *Client) DownloadFile(local string, cloud pkg.File) error {
	var file *os.File
	stat, err := os.Stat(local)
	if err == nil && stat.IsDir() {
		local = path.Clean(local + "/" + cloud.Name())
		file, err = os.OpenFile(local, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(local, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	}
	if err != nil {
		return err
	}
	stat, _ = file.Stat()
	defer file.Close()
	localSize := stat.Size()
	if localSize == cloud.Size() {
		//TODO MD5一致性检查
		return nil
	}
	resp, err := c.api.Download(cloud, localSize)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}
