package web

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
	"github.com/gowsp/cloud189/pkg/file"
)

type searchResp struct {
	Code    int         `json:"res_code,omitempty"`
	Count   int         `json:"count,omitempty"`
	Files   []*fileResp `json:"fileList,omitempty"`
	Folders []*fileResp `json:"folderList,omitempty"`
}

type fileResp struct {
	IsFolder    bool
	ParentId    json.Number
	FileId      json.Number `json:"id,omitempty"`
	FileName    string      `json:"name,omitempty"`
	FileSize    int64       `json:"size,omitempty"`
	MD5         string      `json:"md5,omitempty"`
	FileModTime string      `json:"lastOpTime,omitempty"`
}

func (f *fileResp) Id() string        { return f.FileId.String() }
func (f *fileResp) PId() string       { return f.ParentId.String() }
func (f *fileResp) Name() string      { return f.FileName }
func (f *fileResp) Size() int64       { return f.FileSize }
func (f *fileResp) Mode() os.FileMode { return os.ModePerm }
func (f *fileResp) ModTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", f.FileModTime)
	return t
}
func (f *fileResp) IsDir() bool      { return f.IsFolder }
func (f *fileResp) Sys() interface{} { return nil }
func (f *fileResp) ContentType(ctx context.Context) (string, error) {
	return path.Ext(f.Name()), nil
}
func (f *fileResp) ETag(ctx context.Context) (string, error) {
	return strconv.FormatInt(f.ModTime().UnixMilli(), 10), nil
}

func (c *api) Find(id, name string) (pkg.File, error) {
	if file.IsSystem(id, name) {
		return c.FindDir(id, name)
	}
	return c.FindFile(id, name)
}

func (c *api) FindDir(id, name string) (pkg.File, error) {
	return cache.Find(id, name, func() ([]*folder, error) {
		return c.ListDir(id)
	})
}
func (c *api) FindFile(id, name string) (pkg.File, error) {
	return cache.Find(id, name, func() ([]*fileResp, error) {
		var files searchResp
		err := c.search(&files, id, name, 1)
		return append(files.Folders, files.Files...), err
	})
}

func (c *api) search(result *searchResp, id, name string, page int) error {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	params.Set("filename", name)
	params.Set("recursive", "0")
	params.Set("iconOption", "5")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")
	var files searchResp
	err := c.invoker.Get("/open/file/searchFiles.action", params, &files)
	if err != nil {
		return err
	}
	for _, f := range files.Folders {
		f.IsFolder = true
	}
	parent := json.Number(id)
	for _, f := range files.Files {
		f.ParentId = parent
	}
	if page*100 < files.Count {
		return c.search(result, id, name, page+1)
	}
	return os.ErrNotExist
}
