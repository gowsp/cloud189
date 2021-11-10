package file

import (
	"io"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/gowsp/cloud189-cli/pkg"
)

func IsNetFile(name string) bool {
	r := regexp.MustCompile("^http[s]{0,1}://.*")
	return r.Match([]byte(name))
}

type NetFile struct {
	client   pkg.Client
	parentId string
	url      string
}

func NewNetFile(parentId, url string, client pkg.Client) *NetFile {
	return &NetFile{parentId: parentId, url: url, client: client}
}

func (f *NetFile) Upload() {
	resp, err := http.DefaultClient.Get(f.url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	size := resp.ContentLength
	if size < 1 {
		log.Println("not found content, skip")
		return
	}
	name := path.Base(f.url)
	sf := NewStreamFileWithParent(name, f.parentId, size, f.client)
	io.Copy(sf, resp.Body)
}
