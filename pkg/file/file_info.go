package file

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	KB = 1 << 10
	MB = 1 << 20

	Slice = 10 * MB
)

type FileInfo struct {
	FileId      json.Number `json:"id,omitempty"`
	FileName    string      `json:"name,omitempty"`
	FileSize    int64       `json:"size,omitempty"`
	IsFolder    bool        `json:"isFolder,omitempty"`
	MD5         string      `json:"md5,omitempty"`
	FileModTime string      `json:"lastOpTime,omitempty"`
}

func (f *FileInfo) Id() string        { return f.FileId.String() }
func (f *FileInfo) Name() string      { return f.FileName }
func (f *FileInfo) Size() int64       { return f.FileSize }
func (f *FileInfo) Mode() os.FileMode { return os.ModePerm }
func (f *FileInfo) ModTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", f.FileModTime)
	return t
}
func (f *FileInfo) IsDir() bool      { return f.IsFolder }
func (f *FileInfo) Sys() interface{} { return nil }
func (f *FileInfo) ContentType(ctx context.Context) (string, error) {
	return path.Ext(f.Name()), nil
}
func (f *FileInfo) ETag(ctx context.Context) (string, error) {
	if f.IsDir() {
		return strconv.FormatInt(f.ModTime().Unix(), 10), nil
	}
	return f.MD5, nil
}
