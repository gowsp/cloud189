package webdav

import (
	"context"
	"net/http"
	"os"
	"path"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/drive"
	"golang.org/x/net/webdav"
)

const OPS = "OPS"

type convert interface {
	open(client pkg.App, name string, flag int) (webdav.File, error)
}
type put struct {
	length int64
}

func (p *put) open(client pkg.App, name string, flag int) (webdav.File, error) {
	dir := path.Dir(name)
	parent, err := client.Stat(dir)
	if err != nil {
		return nil, err
	}
	return drive.NewStreamFile(name, parent.Id(), p.length, flag&os.O_TRUNC != 0, client.Uploader()), nil
}

func prepare(req *http.Request) *http.Request {
	switch req.Method {
	case http.MethodPut:
		req = newContext(req, &put{
			length: req.ContentLength,
		})
	}
	return req
}
func newContext(req *http.Request, val convert) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, OPS, val)
	return req.WithContext(ctx)
}
