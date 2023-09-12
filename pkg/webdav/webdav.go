package webdav

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gowsp/cloud189/pkg"
	"golang.org/x/net/webdav"
)

var errInvalidIfHeader = errors.New("webdav: invalid If header")

func Serve(addr string, client pkg.Drive) {
	fs := &CloudFileSystem{
		app: client,
	}
	fs.handler = &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}
	err := http.ListenAndServe(addr, fs)
	if err != nil {
		fmt.Println(err)
	}

}
