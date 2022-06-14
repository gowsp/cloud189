package app

import (
	"net/url"
	"os"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

func (c *api) Rename(src pkg.File, dest string) (err error) {
	if src == nil {
		return os.ErrNotExist
	}
	if src.IsDir() {
		err = c.renameFoler(src.Id(), dest)
	} else {
		err = c.renameFile(src.Id(), dest)
	}
	if err == nil {
		cache.Invalid(src)
	}
	return
}
func (c *api) renameFile(id, dest string) error {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("destFileName", dest)
	var f map[string]interface{}
	return c.invoker.Post("/renameFile.action", params, &f)
}
func (c *api) renameFoler(id, dest string) error {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("destFolderName", dest)
	var result map[string]interface{}
	return c.invoker.Post("/renameFolder.action", params, &result)
}
