package drive

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

func (f *FS) Share(prifix, cloud string) (func(http.ResponseWriter, *http.Request), error) {
	if _, err := f.stat(cloud); err != nil {
		return nil, err
	}
	prifix = strings.TrimRight(prifix, "/")
	return func(w http.ResponseWriter, r *http.Request) {
		target := strings.TrimPrefix(r.RequestURI, prifix)
		target = path.Join(cloud, target)
		log.Println("request", target)
		if f.shareFromCache(target, w) {
			return
		}
		file, err := f.stat(target)
		if err == nil {
			f.doShare(target, file, w)
			return
		}
		if os.IsNotExist(err) {
			w.WriteHeader(404)
			return
		}
		writeError(w, err)
	}, nil
}
func (f *FS) doShare(target string, file pkg.File, w http.ResponseWriter) {
	if file.IsDir() {
		w.WriteHeader(204)
		return
	}
	resp, err := f.api.Download(file, 0)
	if err != nil {
		writeError(w, err)
		return
	}
	defer resp.Body.Close()
	val := resp.Request.URL.String()
	f.share.Store(target, val)
	w.Header().Add("Location", val)
	w.WriteHeader(302)
}
func (f *FS) shareFromCache(target string, w http.ResponseWriter) bool {
	if v, ok := f.share.Load(target); ok {
		u, err := url.Parse(v.(string))
		if err != nil {
			return false
		}
		expires, err := strconv.ParseInt(u.Query().Get("Expires"), 10, 0)
		if err != nil {
			return false
		}
		if expires > time.Now().Unix() {
			w.Header().Add("Location", v.(string))
			w.WriteHeader(302)
			return true
		}
	}
	return false
}
func writeError(w http.ResponseWriter, err error) {
	info := err.Error()
	log.Println("download error", info)
	w.WriteHeader(500)
	w.Write([]byte(info))
}
