package webdav

import (
	"net/http"
	"path/filepath"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

func (h *CloudFileSystem) handlePut(w http.ResponseWriter, r *http.Request) (status int, err error) {
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil {
		return status, err
	}
	if r.ContentLength == 0 {
		return http.StatusCreated, nil
	}
	dir, name := filepath.Split(reqPath)
	parent, err := h.app.Stat(dir)
	if err != nil {
		return http.StatusNotFound, err
	}
	f := file.NewWebFile(parent.(pkg.File).Id(), name, r)
	if copyErr := h.app.UploadFrom(f); copyErr != nil {
		return http.StatusMethodNotAllowed, copyErr
	}
	stat, err := h.app.Stat(reqPath)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("ETag", stat.Name()+stat.ModTime().String())
	return http.StatusCreated, nil
}
