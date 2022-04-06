package web

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

type listResp struct {
	Code  json.Number `json:"res_code,omitempty"`
	Count int         `json:"recordCount,omitempty"`
	Data  []*detail   `json:"data,omitempty"`
}

func (c *api) ListFile(id string) ([]pkg.File, error) {
	return cache.List(id, func() ([]*fileResp, error) { return c.openList(id, 1) })
}

func (c *api) portalList(id string, page int) (result []*detail, err error) {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var file listResp
	err = c.invoker.Get("/portal/listFiles.action", params, file)
	if err != nil {
		return
	}
	result = append(result, file.Data...)
	if 100*page < file.Count {
		var more []*detail
		more, err = c.portalList(id, page+1)
		result = append(result, more...)
	}
	return
}

type listFileResp struct {
	Code json.Number `json:"res_code,omitempty"`
	List *fileList   `json:"fileListAO,omitempty"`
}

type fileList struct {
	Count   int         `json:"count,omitempty"`
	Files   []*fileResp `json:"fileList,omitempty"`
	Folders []*fileResp `json:"folderList,omitempty"`
}

func (l *fileList) files(id string) (data []*fileResp) {
	if l == nil {
		return
	}
	for _, f := range l.Folders {
		f.IsFolder = true
	}
	parentId := json.Number(id)
	for _, f := range l.Files {
		f.ParentId = parentId
	}
	data = append(data, l.Folders...)
	data = append(data, l.Files...)
	return
}

type fileTime time.Time

func (j *fileTime) UnmarshalJSON(b []byte) error {
	json := string(b)
	s := strings.Trim(json, "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = fileTime(t)
	return nil
}

type fileResp struct {
	IsFolder    bool
	ParentId    json.Number
	FileId      json.Number `json:"id,omitempty"`
	FileName    string      `json:"name,omitempty"`
	FileSize    int64       `json:"size,omitempty"`
	MD5         string      `json:"md5,omitempty"`
	FileModTime fileTime    `json:"lastOpTime,omitempty"`
}

func (f *fileResp) Id() string        { return f.FileId.String() }
func (f *fileResp) PId() string       { return f.ParentId.String() }
func (f *fileResp) Name() string      { return f.FileName }
func (f *fileResp) Size() int64       { return f.FileSize }
func (f *fileResp) Mode() os.FileMode { return os.ModePerm }
func (f *fileResp) ModTime() time.Time {
	return time.Time(f.FileModTime)
}
func (f *fileResp) IsDir() bool      { return f.IsFolder }
func (f *fileResp) Sys() interface{} { return nil }
func (f *fileResp) ContentType(ctx context.Context) (string, error) {
	return path.Ext(f.Name()), nil
}
func (f *fileResp) ETag(ctx context.Context) (string, error) {
	return strconv.FormatInt(f.ModTime().Unix(), 10), nil
}

func (c *api) openList(id string, page int) (result []*fileResp, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("mediaType", "0")
	params.Set("orderBy", "lastOpTime")
	params.Set("descending", "true")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var file listFileResp
	err = c.invoker.Get("/open/file/listFiles.action", params, &file)
	if err != nil {
		return
	}
	result = append(result, file.List.files(id)...)
	if 100*page < file.List.Count {
		var more []*fileResp
		more, err = c.openList(id, page+1)
		result = append(result, more...)
	}
	return
}
