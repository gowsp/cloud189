package app

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
	"github.com/gowsp/cloud189/pkg/file"
)

func (c *api) ListFile(id string) ([]pkg.File, error) {
	return cache.List(id, func() ([]*file.FileInfo, error) {
		return c.listFile(file.FileType_All, id, 1)
	})
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

func (c *api) listFile(fileType file.FileType, id string, page int) (result []*file.FileInfo, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("fileType", fileType.String())
	params.Set("mediaType", "0")
	params.Set("mediaAttr", "0")
	params.Set("iconOption", "0")
	params.Set("orderBy", "lastOpTime")
	params.Set("descending", "true")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var resp listFileResp
	err = c.invoker.Get("/listFiles.action", params, &resp)
	if err != nil {
		return
	}
	result = append(result, resp.List.files(id)...)
	if 100*page < resp.List.Count {
		var more []*file.FileInfo
		more, err = c.listFile(fileType, id, page+1)
		result = append(result, more...)
	}
	return
}
