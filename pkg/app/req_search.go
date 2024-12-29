package app

import (
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

func (d *api) Search(parent pkg.File, fileType pkg.FileType, name string) ([]pkg.File, error) {
	return d.search(parent.Id(), strconv.Itoa(int(fileType)), name, 1)
}

type searchResult struct {
	Code    int         `json:"res_code"`
	Message string      `json:"res_message"`
	Count   int         `json:"count"`
	Files   []*fileInfo `json:"fileList"`
	Folders []*folder   `json:"folderList"`
}

func (l *searchResult) fill(id string) (data []pkg.File) {
	if l == nil || l.Count == 0 {
		return
	}
	for _, f := range l.Files {
		f.ParentID = id
		data = append(data, f)
	}
	for _, f := range l.Folders {
		data = append(data, f)
	}
	return
}

func (c *api) search(id, fileType, name string, page int) (result []pkg.File, err error) {
	if file.IsSystem(id, name) {
		return c.List(file.Root, pkg.DIR)
	}
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("filename", name)
	params.Set("fileType", fileType)
	params.Set("mediaType", "0")
	params.Set("mediaAttr", "0")
	params.Set("recursive", "0")
	params.Set("iconOption", "0")
	params.Set("descending", "true")
	params.Set("orderBy", "filename")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	var files searchResult
	err = c.invoker.Get("/searchFiles.action", params, &files)
	if err != nil {
		return
	}
	result = append(result, files.fill(id)...)
	if page*100 < files.Count {
		var more []pkg.File
		more, err = c.search(id, fileType, name, page+1)
		result = append(result, more...)
	}
	return
}
