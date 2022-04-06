package web

import (
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
	"github.com/gowsp/cloud189/pkg/file"
)

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
		return c.search(id, name, 1)
	})
}

func (c *api) search(id, name string, page int) (result []*fileResp, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	params.Set("filename", name)
	params.Set("recursive", "0")
	params.Set("iconOption", "5")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")
	var files fileList
	err = c.invoker.Get("/open/file/searchFiles.action", params, &files)
	if err != nil {
		return
	}
	result = append(result, files.files(id)...)
	if page*100 < files.Count {
		var more []*fileResp
		more, err = c.search(id, name, page+1)
		result = append(result, more...)
	}
	return
}
