package app

import (
	"net/url"
	"strings"

	"github.com/gowsp/cloud189/pkg"
)

func (c *api) Move(target pkg.File, sources ...pkg.File) error {
	if len(sources) == 0 {
		return nil
	}
	list := make([]string, len(sources))
	for i, src := range sources {
		list[i] = src.Id()
	}
	return c.move(target.Id(), list)
}
func (c *api) move(dir string, source []string) error {
	params := make(url.Values)
	params.Set("fileIdList", strings.Join(source, ";"))
	params.Set("destParentFolderId", dir)
	var f map[string]interface{}
	return c.invoker.Post("/batchMoveFile.action", params, &f)
}
