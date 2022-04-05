package web

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

type listResp struct {
	Code  json.Number `json:"res_code,omitempty"`
	Count int         `json:"recordCount,omitempty"`
	Data  []*FileInfo `json:"data,omitempty"`
}

func (c *Api) ListFile(id string) ([]pkg.File, error) {
	return cache.List(id, func(entry *cache.DirEntry) bool {
		ava := c.cached(entry)
		entry.Init = true
		return ava
	}, func() ([]*FileInfo, error) {
		result := &listResp{}
		err := c.list(id, result, 1)
		return result.Data, err
	})
}

func (c *Api) list(id string, result *listResp, page int) error {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")

	file := &listResp{}
	err := c.invoker.Get("/portal/listFiles.action", params, file)
	if err != nil {
		return err
	}
	result.Data = append(result.Data, file.Data...)
	if 100*page < file.Count {
		return c.list(id, result, page+1)
	}
	return nil
}
