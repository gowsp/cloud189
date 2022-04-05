package webdav

import (
	"fmt"
	"net/http"

	"github.com/gowsp/cloud189/pkg"
	"golang.org/x/net/webdav"
)

type WEBDAV_KEY string

const FILE_SIZE WEBDAV_KEY = "FileSize"

func Serve(addr string, client pkg.App) {
	sys := &CloudFileSystem{app: client}
	fs := &webdav.Handler{
		FileSystem: sys,
		LockSystem: webdav.NewMemLS(),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "COPY":
			status, _ := sys.Copy(w, req)
			w.WriteHeader(status)
			if status != http.StatusNoContent {
				w.Write([]byte(webdav.StatusText(status)))
			}
			return
		}
		req = prepare(req)
		fs.ServeHTTP(w, req)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}

}
