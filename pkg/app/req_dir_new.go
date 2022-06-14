package app

import (
	"net/url"
	"path"

	"github.com/gowsp/cloud189/pkg/cache"
)

func (c *api) Mkdir(parentId, path string, parents bool) error {
	if parents {
		_, err := c.Mkdirs(parentId, path)
		return err
	}
	err := c.mkdir(parentId, path)
	if err == nil {
		cache.InvalidId(parentId)
	}
	return err
}
func (c *api) Mkdirs(parentId string, dirs ...string) (map[string]interface{}, error) {
	length := len(dirs)
	if length == 0 {
		return make(map[string]interface{}), nil
	}
	var err error
	for _, v := range dirs {
		if path.IsAbs(v) {
			v = v[1:]
		}
		err = c.mkdir(parentId, v)
	}
	var result map[string]interface{}
	cache.InvalidId(parentId)
	return result, err
}
func (c *api) mkdir(parentId, name string) error {
	var result map[string]interface{}
	dir, base := path.Split(name)
	params := url.Values{"folderName": {base}, "relativePath": {dir}, "parentFolderId": {parentId}}
	return c.invoker.Post("/createFolder.action", params, &result)
}
