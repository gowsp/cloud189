package app

import (
	"log"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

func (c *api) Move(target string, files ...pkg.File) error {
	if len(files) == 0 {
		return nil
	}
	var err error
	for _, src := range files {
		if src.IsDir() {
			err = c.moveFoler(target, src)
		} else {
			err = c.moveFile(target, src)
		}
		if err == nil {
			cache.Invalid(src)
		} else {
			log.Println(err)
		}
	}
	cache.Invalid(files...)
	cache.InvalidId(target)
	return err
}
func (c *api) moveFile(id string, src pkg.File) error {
	params := make(url.Values)
	params.Set("fileId", src.Id())
	params.Set("destFileName", src.Name())
	params.Set("destParentFolderId", id)
	var f map[string]interface{}
	return c.invoker.Post("/moveFile.action", params, &f)
}
func (c *api) moveFoler(id string, src pkg.File) error {
	params := make(url.Values)
	params.Set("folderId", src.Id())
	params.Set("destFolderName", src.Name())
	params.Set("destParentFolderId", id)
	var result map[string]interface{}
	return c.invoker.Post("/moveFolder.action", params, &result)
}
