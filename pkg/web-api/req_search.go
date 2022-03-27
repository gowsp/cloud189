package web

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

type searchResp struct {
	Code    int         `json:"res_code,omitempty"`
	Count   int         `json:"count,omitempty"`
	Files   []*fileResp `json:"fileList,omitempty"`
	Folders []*fileResp `json:"folderList,omitempty"`
}

type fileResp struct {
	ParentId    json.Number
	FileId      json.Number `json:"id,omitempty"`
	FileName    string      `json:"name,omitempty"`
	FileSize    int64       `json:"size,omitempty"`
	IsFolder    bool        `json:"isFolder,omitempty"`
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

func (c *Api) Find(id, name string) (pkg.File, error) {
	return c.searchByName(id, name, 1, func(sr *searchResp) []*fileResp {
		sr.Folders = append(sr.Folders, sr.Files...)
		return sr.Folders
	})
}

func (c *Api) FindDir(id, name string) (pkg.File, error) {
	files, err := c.ListDir(id)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Name() == name {
			return f, nil
		}
	}
	return nil, os.ErrNotExist
}

func (c *Api) FindFile(id, name string) (pkg.File, error) {
	return c.searchByName(id, name, 1, func(sr *searchResp) []*fileResp { return sr.Files })
}
func (c *Api) search(id, name string, page int) (*searchResp, error) {
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
	for _, f := range files.Files {
		f.ParentId = json.Number(id)
	}
	for _, f := range files.Folders {
		f.IsFolder = true
	}
	return &files, err
}

func (c *Api) searchByName(id, name string, page int,
	filter func(*searchResp) []*fileResp) (pkg.File, error) {
	if id == pkg.Root.Id() {
		if val, ok := pkg.System[name]; ok {
			return &pkg.SysFolder{
				FileId:   json.Number(val),
				ParentId: pkg.Root.Id(),
				FileName: name}, nil
		}
	}
	files, err := c.search(id, name, page)
	if err != nil {
		return nil, err
	}
	data := filter(files)
	for _, f := range data {
		if f.Name() == name {
			return f, nil
		}
	}
	if page*100 < files.Count {
		return c.searchByName(id, name, page+1, filter)
	}
	return nil, os.ErrNotExist
}
