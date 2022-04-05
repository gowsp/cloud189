package web

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/gowsp/cloud189/pkg/cache"
)

type mkdirs struct {
	ParentId string   `json:"parentId,omitempty"`
	Paths    []string `json:"paths,omitempty"`
}
type response struct {
	Code int    `json:"res_code,omitempty"`
	Msg  string `json:"res_message,omitempty"`
}

func (c *api) Mkdir(parentId, path string, parents bool) error {
	if parents {
		_, err := c.Mkdirs(parentId, path)
		return err
	}
	return c.mkdir(parentId, path)
}
func (c *api) Mkdirs(parentId string, dirs ...string) (map[string]interface{}, error) {
	length := len(dirs)
	if length == 0 {
		return make(map[string]interface{}), nil
	}
	paths := make([]string, length)
	for i, v := range dirs {
		if path.IsAbs(v) {
			paths[i] = v[1:]
		} else {
			paths[i] = v
		}
	}
	data := mkdirs{ParentId: parentId, Paths: paths}
	val, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	params := url.Values{"folderList": {string(val)}, "opScene": {"1"}}
	err = c.invoker.Post("/portal/createFolders.action", params, &result)
	if err == nil && result["res_code"].(float64) == 0 {
		cache.InvalidId(parentId)
		return result, nil
	}
	return result, err
}
func (c *api) mkdir(parentId, name string) error {
	var result map[string]interface{}
	params := url.Values{"folderName": {name}, "parentFolderId": {parentId}}
	err := c.invoker.Post("/open/file/createFolder.action", params, &result)
	if err == nil && result["res_code"].(float64) == 0 {
		cache.InvalidId(parentId)
		return nil
	}
	return err
}
