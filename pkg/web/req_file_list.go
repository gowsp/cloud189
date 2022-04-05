package web

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

type listResp struct {
	Code  json.Number `json:"res_code,omitempty"`
	Count int         `json:"recordCount,omitempty"`
	Data  []*detail   `json:"data,omitempty"`
}

func (c *api) ListFile(id string) ([]pkg.File, error) {
	return cache.List(id, func() ([]*fileResp, error) { return c.openList(id) })
}

func (c *api) portalList(id string) ([]*detail, error) {
	var result listResp
	err := c.portalListFile(id, &result, 1)
	return result.Data, err
}
func (c *api) portalListFile(id string, result *listResp, page int) error {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	file := &listResp{}
	err := c.invoker.Get("/portal/listFiles.action", params, file)
	if err != nil {
		return err
	}
	result.Data = append(result.Data, file.Data...)
	if 100*page < file.Count {
		return c.portalListFile(id, result, page+1)
	}
	return nil
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

func (c *api) openList(id string) ([]*fileResp, error) {
	var result fileList
	err := c.openListFile(&result, id, 1)
	return append(result.Folders, result.Files...), err
}
func (c *api) openListFile(result *fileList, id string, page int) error {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("mediaType", "0")
	params.Set("orderBy", "lastOpTime")
	params.Set("descending", "true")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	var file listFileResp
	err := c.invoker.Get("/open/file/listFiles.action", params, &file)
	if err != nil {
		return err
	}
	for _, f := range file.List.Folders {
		f.IsFolder = true
	}
	parentId := json.Number(id)
	for _, f := range file.List.Files {
		f.ParentId = parentId
	}
	result.Files = append(result.Files, file.List.Files...)
	result.Folders = append(result.Folders, file.List.Folders...)
	if 100*page < file.List.Count {
		return c.openListFile(result, id, page+1)
	}
	return nil
}
