package webdav

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gowsp/cloud189-cli/pkg"
	"golang.org/x/net/webdav"
)

type WEBDAV_KEY string

const FILE_SIZE WEBDAV_KEY = "FileSize"

func Serve(addr string, client pkg.Client) {
	fs := &webdav.Handler{
		FileSystem: &Cloud189FileSystem{client: client},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, e error) {
			data, _ := httputil.DumpRequest(r, true)
			log.Println(string(data))
		},
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
