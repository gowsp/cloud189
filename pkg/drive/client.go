package drive

import (
	"os"
	"strings"

	"github.com/gowsp/cloud189/pkg"
)

type Client struct {
	api pkg.Api
}

func NewClient(api pkg.Api) *Client {
	return &Client{api: api}
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
	var file pkg.File = pkg.Root
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
func (f *Client) Mkdir(path string, parents bool) error {
	if parents {
		return f.api.Mkdir(pkg.Root.Id(), path, parents)
	}
	dir, err := f.Stat(Dir(path))
	if err != nil {
		return err
	}
	return f.api.Mkdir(dir.Id(), Base(path), parents)
}
func (f *Client) Mkdirs(path ...string) error {
	return f.api.Mkdirs(pkg.Root.Id(), path...)
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
	if err != nil {
		return err
	}
	src := f.parse(from...)
	if len(src) == 0 {
		return nil
	}
	if dest.IsDir() {
		return f.api.Copy(dest.Id(), src...)
	}
	return os.ErrInvalid
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
	if dest.IsDir() {
		return f.api.Copy(dest.Id(), src)
	}
	return os.ErrInvalid
}

func (f *Client) Move(target string, src ...string) error {
	size := len(src)
	if size == 0 {
		return nil
	}
	tar, err := f.Stat(target)
	if err != nil {
		return err
	}
	if len(src) == 1 {
		return f.move(src[0], target)
	}
	if !tar.IsDir() {
		return os.ErrInvalid
	}
	data := f.parse(src...)
	if len(data) == 0 {
		return nil
	}
	return f.api.Move(tar.Id(), data...)
}
func (f *Client) move(oldName, newName string) error {
	if oldName == newName {
		return nil
	}
	src, err := f.Stat(oldName)
	if err != nil {
		return err
	}
	dest, err := f.Stat(newName)
	if os.IsNotExist(err) {
		if IsDir(newName) {
			return err
		}
		dir := Dir(newName)
		if dir == Dir(oldName) {
			return f.api.Rename(src, Base(newName))
		}
		dest, err = f.Stat(dir)
		if err != err {
			return err
		}
		if Base(newName) == Base(oldName) {
			return f.api.Move(dest.PId(), src)
		} else {
			f.api.Rename(src, Base(newName))
			return f.api.Move(dest.PId(), src)
		}
	}
	if dest.IsDir() {
		return f.api.Move(dest.Id(), src)
	}
	err = f.api.Delete(dest)
	if err != nil {
		return err
	}
	if Base(newName) == Base(oldName) {
		return f.api.Move(dest.PId(), src)
	} else {
		f.api.Rename(src, Base(newName))
		return f.api.Move(dest.PId(), src)
	}
}
