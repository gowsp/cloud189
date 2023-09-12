package app

import (
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
)

func (d *api) List(parent pkg.File, fileType pkg.FileType) ([]pkg.File, error) {
	id := parent.Id()
	return d.list(id, strconv.Itoa(int(fileType)), 1)
}

type listFileResp struct {
	Code    int    `json:"res_code"`
	Message string `json:"res_message"`
	Result  struct {
		Count int `json:"count"`
		Size  int `json:"fileListSize"`

		Files   []*fileInfo `json:"fileList"`
		Folders []*folder   `json:"folderList"`
	} `json:"fileListAO"`
	LastRev int64 `json:"lastRev"`
}

func (l *listFileResp) fill(id string) (data []pkg.File) {
	if l == nil || l.Result.Count == 0 {
		return
	}
	for _, f := range l.Result.Files {
		f.ParentID = id
		data = append(data, f)
	}
	for _, f := range l.Result.Folders {
		data = append(data, f)
	}
	return
}

func (c *api) list(id, fileType string, page int) (result []pkg.File, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("fileType", fileType)
	params.Set("mediaType", "0")
	params.Set("mediaAttr", "0")
	params.Set("iconOption", "0")
	params.Set("orderBy", "filename")
	params.Set("descending", "true")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var resp listFileResp
	err = c.invoker.Get("/listFiles.action", params, &resp)
	if err != nil {
		return
	}
	result = append(result, resp.fill(id)...)
	if 100*page < resp.Result.Count {
		var more []pkg.File
		more, err = c.list(id, fileType, page+1)
		result = append(result, more...)
	}
	return
}
