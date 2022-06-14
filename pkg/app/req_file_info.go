package app

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
)

func (c *api) Detail(id string) (string, error) {
	var info map[string]string
	err := c.invoker.Get("/getFileDownloadUrl.action", url.Values{"fileId": {id}}, &info)
	return info["fileDownloadUrl"], err
}

func (c *api) Download(file pkg.File, start int64) (*http.Response, error) {
	if file.IsDir() {
		return nil, errors.New("not support download dir")
	}
	url, _ := c.Detail(file.Id())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, file.Size()))
	return c.invoker.Send(req)
}
