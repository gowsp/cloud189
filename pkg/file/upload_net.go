package file

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/gowsp/cloud189/pkg"
)

func IsNetFile(name string) bool {
	r := regexp.MustCompile("^http[s]{0,1}://.*")
	return r.Match([]byte(name))
}
func NewURLFile(parentId, link string) pkg.Upload {
	u, _ := url.Parse(link)
	resp, _ := http.Get(link)
	size := resp.ContentLength
	pr, pw := io.Pipe()
	num := int(math.Ceil(float64(size) / float64(Slice)))
	return &NetFile{
		parentId: parentId,
		name:     path.Base(u.Path),
		size:     size,
		pr:       pr,
		pw:       pw,
		num:      int64(num),
		md5:      md5.New(),
		slices:   make([]string, num),
		data:     resp.Body,
	}
}
func NewWebFile(parentId, name string, data *http.Request) pkg.Upload {
	size := data.ContentLength
	pr, pw := io.Pipe()
	num := int(math.Ceil(float64(size) / float64(Slice)))
	return &NetFile{
		parentId:  parentId,
		name:      name,
		size:      size,
		pr:        pr,
		pw:        pw,
		num:       int64(num),
		md5:       md5.New(),
		slices:    make([]string, num),
		data:      data.Body,
		overwrite: true,
	}
}

type NetFile struct {
	parentId  string
	name      string
	size      int64
	num       int64
	overwrite bool
	slices    []string
	md5       hash.Hash
	pr        *io.PipeReader
	pw        *io.PipeWriter
	data      io.ReadCloser
	start     sync.Once
	close     bool
}

func (f *NetFile) ParentId() string {
	return f.parentId
}
func (f *NetFile) Name() string {
	return f.name
}
func (f *NetFile) Overwrite() bool {
	return f.overwrite
}
func (f *NetFile) Size() int64 {
	return f.size
}
func (f *NetFile) SliceNum() int {
	return len(f.slices)
}
func (f *NetFile) LazyCheck() bool {
	return true
}
func (f *NetFile) FileMD5() string {
	return hex.EncodeToString(f.md5.Sum(nil))
}
func (f *NetFile) SliceMD5() string {
	if f.num == 1 {
		return f.FileMD5()
	}
	m := md5.New()
	m.Write([]byte(strings.Join(f.slices, "\n")))
	return hex.EncodeToString(m.Sum(nil))
}
func (f *NetFile) Part(i int64) pkg.UploadPart {
	f.start.Do(f.copy)
	m := md5.New()
	buff := bytes.NewBuffer(nil)
	io.Copy(io.MultiWriter(m, buff), io.LimitReader(f.pr, Slice))
	v := m.Sum(nil)
	f.slices[i] = strings.ToUpper(hex.EncodeToString(v))
	name := base64.StdEncoding.EncodeToString(v)
	return &FilePart{data: buff, name: name, num: i}
}

func (f *NetFile) copy() {
	go func() {
		w := io.MultiWriter(f.md5, f.pw)
		io.Copy(w, f.data)
		f.close = true
		f.data.Close()
		f.pw.Close()
	}()
}
