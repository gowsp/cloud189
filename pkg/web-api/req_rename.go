package web

import (
	"net/url"

	"github.com/gowsp/cloud189/pkg"
)

func (c *Api) Rename(file pkg.File, dest string) error {
	if file.IsDir() {
		return c.renameFoler(file.Id(), dest)
	} else {
		return c.renameFile(file.Id(), dest)
	}
}
func (c *Api) renameFile(id, dest string) error {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("destFileName", dest)
	var f FileInfo
	return c.invoker.Post("/open/file/renameFile.action", params, &f)
}
func (c *Api) renameFoler(id, dest string) error {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("destFolderName", dest)
	return c.invoker.Post("/open/folder/renameFolder.action", params, nil)
}
