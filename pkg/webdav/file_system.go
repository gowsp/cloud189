package webdav

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/net/webdav"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

type Cloud189FileSystem struct {
	client pkg.Client
}

func (f *Cloud189FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return f.client.Mkdir(name)
}

func (f *Cloud189FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// if flag&os.O_CREATE != 0 {
	// 	size := ctx.Value(FILE_SIZE).(int64)
	// 	return file.NewStreamFile(name, size, f.client), nil
	// }
	stat, err := f.client.Stat(name)
	if err != nil {
		return nil, err
	}
	return &file.ReadableFile{
		Client:   f.client,
		Id:       json.Number(stat.Id()),
		Name:     stat.Name(),
		IsFolder: stat.IsDir(),
		FileInfo: stat,
	}, nil
}
func (f *Cloud189FileSystem) RemoveAll(ctx context.Context, name string) error {
	return f.client.RemoveAll(name)
}
func (f *Cloud189FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return f.client.Rename(oldName, newName)
}
func (f *Cloud189FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return f.client.Stat(name)
}
