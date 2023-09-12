package drive

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

func (f *FS) Move(target string, source ...string) error {
	if len(source) == 1 {
		return f.singleMove(target, source[0])
	}
	return f.multiMove(target, source...)
}

func (f *FS) singleMove(target string, sources string) error {
	files := f.resolve(sources)
	if len(files) == 0 {
		return fs.ErrNotExist
	}
	source := files[0]
	dest, err := f.stat(target)
	defer func() {
		invalid(source, dest)
	}()
	if err == nil {
		if dest.IsDir() {
			return f.api.Move(dest, files...)
		}
		f.api.Delete(dest)
		if source.PId() == dest.PId() {
			return f.api.Rename(source, dest.Name())
		} else {
			f.api.Move(dest, files...)
			return f.api.Rename(source, dest.Name())
		}
	}
	if errors.Is(err, fs.ErrNotExist) {
		dir, name := filepath.Split(target)
		parent, err := f.stat(dir)
		if err != nil {
			return err
		}
		if err := f.api.Move(parent, source); err != nil {
			return err
		}
		if source.Name() == name {
			return nil
		}
		return f.api.Rename(source, name)
	}
	return err
}
func (f *FS) multiMove(target string, source ...string) error {
	dest, err := f.stat(target)
	if err != nil {
		return err
	}
	if !dest.IsDir() {
		return fmt.Errorf("target '%s' is not a directory", target)
	}
	files := f.resolve(source...)
	defer func() {
		load(dest.Id()).invalid()
		invalid(files...)
	}()
	if len(files) == 0 {
		return nil
	}
	return f.api.Move(dest, files...)
}
