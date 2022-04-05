package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

type detail struct {
	ParentId    json.Number `json:"parentId,omitempty"`
	FileId      json.Number `json:"fileId,omitempty"`
	FileName    string      `json:"fileName,omitempty"`
	FileSize    int64       `json:"fileSize,omitempty"`
	IsFolder    bool        `json:"isFolder,omitempty"`
	FileModTime int64       `json:"lastOpTime,omitempty"`
	CreateTime  int64       `json:"createTime,omitempty"`
	FileCount   int64       `json:"subFileCount,omitempty"`
	DownloadUrl string      `json:"downloadUrl,omitempty"`
}

func (f *detail) Id() string         { return f.FileId.String() }
func (f *detail) PId() string        { return f.ParentId.String() }
func (f *detail) Name() string       { return f.FileName }
func (f *detail) Size() int64        { return f.FileSize }
func (f *detail) Mode() os.FileMode  { return os.ModePerm }
func (f *detail) ModTime() time.Time { return time.UnixMilli(f.FileModTime) }
func (f *detail) IsDir() bool        { return f.IsFolder }
func (f *detail) Sys() interface{} {
	return pkg.FileExt{
		FileCount:   f.FileCount,
		DownloadUrl: "https:" + f.DownloadUrl,
		CreateTime:  time.UnixMilli(f.CreateTime),
	}
}
func (f *detail) ContentType(ctx context.Context) (string, error) {
	return path.Ext(f.Name()), nil
}
func (f *detail) ETag(ctx context.Context) (string, error) {
	return strconv.FormatInt(f.FileModTime, 10), nil
}

func (c *api) Detail(id string) (pkg.File, error) {
	var info detail
	err := c.invoker.Get("/portal/getFileInfo.action", url.Values{"fileId": {id}}, &info)
	return &info, err
}

func (c *api) Download(file pkg.File, start int64) (*http.Response, error) {
	if file.IsDir() {
		return nil, errors.New("not support download dir")
	}
	file, _ = c.Detail(file.Id())
	req, err := http.NewRequest(http.MethodGet, file.Sys().(pkg.FileExt).DownloadUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, file.Size()))
	return c.invoker.http.Do(req)
}
