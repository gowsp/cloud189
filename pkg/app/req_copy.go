package app

import (
	"log"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

func (c *api) Copy(target pkg.File, files ...pkg.File) error {
	var err error
	for _, v := range files {
		err = c.copy(target, v)
		if err != nil {
			log.Println(err)
		}
	}
	cache.InvalidId(target.Id())
	return err
}

func (c *api) copy(targetFolderId pkg.File, src pkg.File) error {
	params := make(url.Values)
	params.Set("fileId", src.Id())
	params.Set("destFileName", src.Name())
	params.Set("destParentFolderId", targetFolderId.Id())
	var result map[string]interface{}
	return c.invoker.Post("/copyFile.action", params, &result)
}
