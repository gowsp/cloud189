package drive

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

type Client struct {
	api pkg.Api
}

func NewClient(api pkg.Api) *Client {
	return &Client{api: api}
}

func (f *Client) Uploader() pkg.Uploader {
	return f.api
}
func (f *Client) Login(name, password string) error {
	return f.api.Login(name, password)
}
func (f *Client) Sign() error {
	return f.api.Sign()
}
func (f *Client) Space() (pkg.Space, error) {
	return f.api.Space()
}
func (f *Client) Stat(name string) (pkg.File, error) {
	var err error
	var file pkg.File = file.Root
	path := strings.Split(name, "/")
	size := len(path) - 1
	for i := 1; i < size; i++ {
		file, err = f.api.FindDir(file.Id(), path[i])
		if err != nil {
			return nil, err
		}
	}
	if path[size] == "" {
		return file, nil
	}
	return f.api.Find(file.Id(), path[size])
}
func (f *Client) List(file pkg.File) ([]pkg.File, error) {
	if file.IsDir() {
		return f.api.ListFile(file.Id())
	}
	return nil, os.ErrInvalid
}
func (f *Client) ListDir(name string) ([]pkg.File, error) {
	stat, err := f.Stat(name)
	if err != nil {
		return nil, err
	}
	return f.api.ListDir(stat.Id())
}
func (f *Client) Mkdir(name string, parents bool) error {
	if parents {
		return f.api.Mkdir(file.Root.Id(), name, parents)
	}
	stat, err := f.Stat(name)
	if err == nil && stat != nil {
		return os.ErrExist
	}
	dir, file := path.Split(name)
	parent, err := f.Stat(dir)
	if err != nil {
		return err
	}
	return f.api.Mkdir(parent.Id(), file, parents)
}
func (f *Client) Mkdirs(path ...string) error {
	_, err := f.api.Mkdirs(file.Root.Id(), path...)
	return err
}
func (f *Client) Remove(paths ...string) error {
	data := f.parse(paths...)
	if len(data) == 0 {
		return nil
	}
	return f.api.Delete(data...)
}
func (f *Client) Copy(target string, from ...string) error {
	size := len(from)
	if size == 0 {
		return nil
	}
	if size == 1 {
		return f.copy(target, from[0])
	}
	dest, err := f.Stat(target)
	if err != nil || !dest.IsDir() {
		return fmt.Errorf("%s: file does not exist or not a directory", target)
	}
	src := f.parse(from...)
	if len(src) == 0 {
		return nil
	}
	return f.api.Copy(dest.Id(), src...)
}
func (f *Client) copy(target string, from string) error {
	if from == target {
		return nil
	}
	src, err := f.Stat(from)
	if err != nil {
		return err
	}
	dest, err := f.Stat(target)
	if err == nil && dest.IsDir() {
		return f.api.Copy(dest.Id(), src)
	}
	if !os.IsNotExist(err) {
		return err
	}
	odir := path.Dir(from)
	ndir := path.Dir(target)
	if ndir == odir {
		return fmt.Errorf("same dir not support copy file")
	}
	dest, err = f.Stat(ndir)
	if err != nil {
		return err
	}
	return f.api.Copy(dest.Id(), src)
}

func (f *Client) Move(target string, src ...string) error {
	size := len(src)
	if size == 0 {
		return nil
	}
	if size == 1 {
		return f.move(src[0], target)
	}
	dest, err := f.Stat(target)
	if err != nil || !dest.IsDir() {
		return fmt.Errorf("%s: file does not exist or not a directory", target)
	}
	files := f.parse(src...)
	if len(files) == 0 {
		return nil
	}
	return f.api.Move(dest.Id(), files...)
}
func (f *Client) move(oldName, newName string) error {
	if oldName == newName {
		return nil
	}
	src, err := f.Stat(oldName)
	if err != nil {
		return fmt.Errorf("%s: file does not exist", oldName)
	}
	ndir, nname := path.Split(newName)
	odir, oname := path.Split(oldName)
	if odir == ndir {
		// same dir rename file
		return f.api.Rename(src, nname)
	}
	dest, err := f.Stat(newName)
	if os.IsNotExist(err) {
		dest, err = f.Stat(ndir)
		if err != err {
			return err
		}
		if nname != oname {
			f.api.Rename(src, nname)
		}
		return f.api.Move(dest.Id(), src)
	}
	if dest.IsDir() {
		return f.api.Move(dest.Id(), src)
	}
	return os.ErrExist
}
