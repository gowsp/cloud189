package drive

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gowsp/cloud189/pkg"
)

func IsNetFile(name string) bool {
	r := regexp.MustCompile("^http[s]{0,1}://.*")
	return r.Match([]byte(name))
}

type NetFile struct {
	client   pkg.Uploader
	parentId string
	url      string
}

func NewNetFile(parentId, url string, client pkg.Uploader) *NetFile {
	return &NetFile{parentId: parentId, url: url, client: client}
}

func (f *NetFile) Upload() {
	resp, err := http.DefaultClient.Get(f.url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	size := resp.ContentLength
	if size < 1 {
		fmt.Println("not found content, skip")
		return
	}
	var name string
	header := resp.Header.Get("Content-Disposition")
	if header == "" {
		name = path.Base(f.url)
	} else {
		name = strings.Split(header, "filename=")[1]
		name, _ = url.PathUnescape(name)
	}
	sf := NewStreamFileWithParent(name, f.parentId, size, f.client)
	io.Copy(sf, resp.Body)
	resp.Body.Close()
}
