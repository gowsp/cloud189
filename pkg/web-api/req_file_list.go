package web

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
)

type listResp struct {
	Code  json.Number `json:"res_code,omitempty"`
	Count int         `json:"recordCount,omitempty"`
	Data  []*FileInfo `json:"data,omitempty"`
}

func (c *Api) ListFile(id string) ([]pkg.File, error) {
	result := &listResp{}
	err := c.list(id, result, 1)
	if err != nil {
		return nil, err
	}
	data := make([]pkg.File, 0)
	for _, f := range result.Data {
		data = append(data, f)
	}
	return data, nil
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
