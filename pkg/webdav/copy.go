package webdav

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var errPrefixMismatch = errors.New("webdav: prefix mismatch")
var errDestinationEqualsSource = errors.New("webdav: destination equals source")
var errInvalidDestination = errors.New("webdav: invalid destination")

func (h *CloudFileSystem) stripPrefix(p string) (string, int, error) {
	if h.Prefix == "" {
		return p, http.StatusOK, nil
	}
	if r := strings.TrimPrefix(p, h.Prefix); len(r) < len(p) {
		return r, http.StatusOK, nil
	}
	return p, http.StatusNotFound, errPrefixMismatch
}

func (h *CloudFileSystem) Copy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	hdr := r.Header.Get("Destination")
	if hdr == "" {
		return http.StatusBadRequest, errInvalidDestination
	}
	u, err := url.Parse(hdr)
	if err != nil {
		return http.StatusBadRequest, errInvalidDestination
	}
	if u.Host != "" && u.Host != r.Host {
		return http.StatusBadGateway, errInvalidDestination
	}

	src, status, err := h.stripPrefix(r.URL.Path)
	if err != nil {
		return status, err
	}

	dst, status, err := h.stripPrefix(u.Path)
	if err != nil {
		return status, err
	}

	if dst == "" {
		return http.StatusBadGateway, errInvalidDestination
	}
	if dst == src {
		return http.StatusForbidden, errDestinationEqualsSource
	}
	err = h.app.Copy(dst, src)
	if err != nil {
		return http.StatusForbidden, err
	}
	return http.StatusCreated, nil
}
