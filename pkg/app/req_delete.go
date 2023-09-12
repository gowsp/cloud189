package app

import (
	"net/url"
	"strings"

	"github.com/gowsp/cloud189/pkg"
)

func (c *api) Delete(files ...pkg.File) error {
	if len(files) == 0 {
		return nil
	}
	list := make([]string, len(files))
	var err error
	for i, src := range files {
		list[i] = src.Id()
	}
	c.deleteFile(list)
	return err
}
func (c *api) deleteFile(list []string) error {
	params := make(url.Values)
	id := strings.Join(list, ";")
	params.Set("fileIdList", id)
	var f map[string]interface{}
	return c.invoker.Post("/batchDeleteFile.action", params, &f)
}
