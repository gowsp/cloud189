package webdav

import (
	"context"
	"net/http"

	"github.com/gowsp/cloud189-cli/pkg"
	"golang.org/x/net/webdav"
)

type WEBDAV_KEY string

const FILE_SIZE WEBDAV_KEY = "FileSize"

func Serve(addr string, client pkg.Client) {
	fs := &webdav.Handler{
		FileSystem: &Cloud189FileSystem{client: client},
		LockSystem: webdav.NewMemLS(),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			ctx := req.Context()
			ctx = context.WithValue(ctx, FILE_SIZE, req.ContentLength)
			req = req.WithContext(ctx)
		}
		fs.ServeHTTP(w, req)
	})
	http.ListenAndServe(addr, nil)

}
