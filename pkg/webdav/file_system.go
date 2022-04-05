package webdav

import (
	"context"
	"os"

	"golang.org/x/net/webdav"

	"github.com/gowsp/cloud189/pkg"
)

type CloudFileSystem struct {
	app    pkg.App
	Prefix string
}

func (f *CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return f.app.Mkdir(name, false)
}

func (f *CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if flag&os.O_CREATE != 0 {
		return ctx.Value(OPS).(convert).open(f.app, name, flag)
	}
	return newRead(f.app, name)
}
func (f *CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	return f.app.Remove(name)
}
func (f *CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return f.app.Move(newName, oldName)
}
func (f *CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return f.app.Stat(name)
}
