package app

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
	if file.IsSystem(id, name) {
		return cache.Find(id, name, func() ([]*file.FileInfo, error) {
			return c.listFile(file.FileType_Dir, id, 1)
		})
	}
	return cache.Find(id, name, func() ([]*file.FileInfo, error) {
		return c.search(file.FileType_Dir, id, name, 1)
	})
}
func (c *api) FindFile(id, name string) (pkg.File, error) {
	return cache.Find(id, name, func() ([]*file.FileInfo, error) {
		return c.search(file.FileType_All, id, name, 1)
	})
}

func (c *api) search(fileType file.FileType, id, name string, page int) (result []*file.FileInfo, err error) {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("filename", name)
	params.Set("fileType", fileType.String())
	params.Set("mediaType", "0")
	params.Set("mediaAttr", "0")
	params.Set("recursive", "0")
	params.Set("iconOption", "0")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	var files fileList
	err = c.invoker.Get("/searchFiles.action", params, &files)
	if err != nil {
		return
	}
	result = append(result, files.files(id)...)
	if page*100 < files.Count {
		var more []*file.FileInfo
		more, err = c.search(fileType, id, name, page+1)
		result = append(result, more...)
	}
	return
}
