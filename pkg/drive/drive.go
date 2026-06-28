package drive

import (
	"errors"
	"fmt"
	"io/fs"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

func New(api pkg.DriveApi) pkg.Drive {
	return &FS{api: api, root: file.Root}
}

type FS struct {
	root  pkg.File
	api   pkg.DriveApi
	share sync.Map
}

func (f *FS) Login(username, password string) error {
	return f.api.PwdLogin(username, password)
}
func (f *FS) QrLogin() error {
	return f.api.QrLogin()
}
func (f *FS) Space() (pkg.Space, error) {
	return f.api.Space()
}

func (f *FS) Open(name string) (fs.File, error) {
	info, err := f.stat(name)
	if err != nil {
		return nil, err
	}
	return f.NewFile(info), nil
}

func (f *FS) Mkdir(name string) error {
	if len(name) == 0 {
		return nil
	}
	_, err := f.stat(name)
	if err == nil {
		return fs.ErrExist
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	dir, err := f.api.Mkdir(file.Root, name)
	if err != nil {
		return err
	}
	invalid(dir)
	return nil
}

func (f *FS) Copy(target string, source ...string) error {
	dest, err := f.stat(target)
	if err != nil || !dest.IsDir() {
		return fmt.Errorf("%s: file does not exist or not a directory", target)
	}
	src := f.resolve(source...)
	if len(src) == 0 {
		return nil
	}
	defer func() {
		load(dest.Id()).invalid()
		invalid(src...)
	}()
	return f.api.Copy(dest, src...)
}

func (f *FS) Delete(name ...string) error {
	files := f.resolve(name...)
	if len(files) == 0 {
		return nil
	}
	err := f.api.Delete(files...)
	for _, file := range files {
		load(file.PId()).delete(file)
	}
	invalid(files...)
	return err
}

func (f *FS) Usage(name string) (pkg.Usage, error) {
	fileInfo, err := f.stat(name)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return f.api.DirUsage(fileInfo)
	}
	return file.NewFileUsage(fileInfo), nil
}
