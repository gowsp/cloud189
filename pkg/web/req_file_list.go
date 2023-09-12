package web

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
	"github.com/gowsp/cloud189/pkg/file"
)

type listResp struct {
	Code  json.Number `json:"res_code,omitempty"`
	Count int         `json:"recordCount,omitempty"`
	Data  []*detail   `json:"data,omitempty"`
}

func (c *api) ListFile(id string) ([]pkg.File, error) {
	return cache.List(id, func() ([]*file.FileInfo, error) { return c.openList(id, 1) })
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
	Count   int              `json:"count,omitempty"`
	Files   []*file.FileInfo `json:"fileList,omitempty"`
	Folders []*file.FileInfo `json:"folderList,omitempty"`
}

func (l *fileList) files(id string) (data []*file.FileInfo) {
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

func (c *api) openList(id string, page int) (result []*file.FileInfo, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("mediaType", "0")
	params.Set("orderBy", "lastOpTime")
	params.Set("descending", "true")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var resp listFileResp
	err = c.invoker.Get("/open/file/listFiles.action", params, &resp)
	if err != nil {
		return
	}
	result = append(result, resp.List.files(id)...)
	if 100*page < resp.List.Count {
		var more []*file.FileInfo
		more, err = c.openList(id, page+1)
		result = append(result, more...)
	}
	return
}
