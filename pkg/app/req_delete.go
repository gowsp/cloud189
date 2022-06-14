package app

import (
	"log"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

func (c *api) Delete(files ...pkg.File) error {
	if len(files) == 0 {
		return nil
	}
	var err error
	for _, src := range files {
		if src.IsDir() {
			err = c.deleteFoler(src.Id())
		} else {
			err = c.deleteFile(src.Id())
		}
		if err == nil {
			cache.Invalid(src)
		} else {
			log.Println(err)
		}
	}
	cache.Delete(files...)
	return err
}
func (c *api) deleteFile(id string) error {
	params := make(url.Values)
	params.Set("fileId", id)
	var f map[string]interface{}
	return c.invoker.Post("/deleteFile.action", params, &f)
}
func (c *api) deleteFoler(id string) error {
	params := make(url.Values)
	params.Set("folderId", id)
	var result map[string]interface{}
	return c.invoker.Post("/deleteFolder.action", params, &result)
}
