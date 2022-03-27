package web

import (
	"encoding/json"
	"net/url"
)

type mkdirs struct {
	ParentId string   `json:"parentId,omitempty"`
	Paths    []string `json:"paths,omitempty"`
}
type response struct {
	Code int    `json:"res_code,omitempty"`
	Msg  string `json:"res_message,omitempty"`
}

func (c *Api) Mkdir(parentId, path string, parents bool) error {
	if parents {
		return c.Mkdirs(parentId, path)
	}
	return c.mkdir(parentId, path)
}
func (c *Api) Mkdirs(parentId string, path ...string) error {
	data := mkdirs{ParentId: parentId, Paths: path}
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var result response
	params := url.Values{"folderList": {string(val)}, "opScene": {"1"}}
	return c.invoker.Post("/portal/createFolders.action", params, &result)
}
func (c *Api) mkdir(parentId, name string) error {
	var file FileInfo
	params := url.Values{"folderName": {name}, "parentFolderId": {parentId}}
	return c.invoker.Post("/open/file/createFolder.action", params, &file)
}
