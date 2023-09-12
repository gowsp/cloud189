package webdav

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"

	"github.com/gowsp/cloud189/pkg"
)

var errUnsupportedMethod = errors.New("webdav: unsupported method")

type CloudFileSystem struct {
	app     pkg.Drive
	Prefix  string
	handler *webdav.Handler
}

func (h *CloudFileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	useNative := true
	status, err := http.StatusBadRequest, errUnsupportedMethod
	switch r.Method {
	case "PUT":
		useNative = false
		status, err = h.handlePut(w, r)
	case "COPY":
		useNative = false
		status, err = h.handleCopyMove(w, r)
	}
	if useNative {
		h.handler.ServeHTTP(w, r)
		return
	}
	if err != nil {
		log.Println(err)
	}
	if status != 0 {
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			w.Write([]byte(webdav.StatusText(status)))
		}
	}
}

func (f *CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return f.app.Mkdir(name)
}
func (f *CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Println("open file", name)
	if flag&os.O_CREATE != 0 {
		return empty, nil
	}
	return newRead(f.app, name)
}
func (f *CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	return f.app.Delete(name)
}
func (f *CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return f.app.Move(newName, oldName)
}
func (f *CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return f.app.Stat(name)
}
