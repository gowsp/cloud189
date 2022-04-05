package web

import (
	"net/url"
	"os"

	"github.com/gowsp/cloud189/pkg"
)

func (c *Api) Rename(src pkg.File, dest string) error {
	if src == nil {
		return os.ErrNotExist
	}
	if src.IsDir() {
		return c.renameFoler(src.Id(), dest)
	} else {
		return c.renameFile(src.Id(), dest)
	}
}
func (c *Api) renameFile(id, dest string) error {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("destFileName", dest)
	var f map[string]interface{}
	return c.invoker.Post("/open/file/renameFile.action", params, &f)
}
func (c *Api) renameFoler(id, dest string) error {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("destFolderName", dest)
	var result map[string]interface{}
	return c.invoker.Post("/open/file/renameFolder.action", params, &result)
}
