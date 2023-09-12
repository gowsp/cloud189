package file

import (
	"regexp"
	"strconv"

	"github.com/gowsp/cloud189/pkg"
)

func IsFastFile(name string) bool {
	r := regexp.MustCompile(`^fast://(\w+):(\d+)/(.+)`)
	return r.Match([]byte(name))
}

type FastFile struct {
	parentId  string
	name      string
	size      int64
	fileMd5   string
	overwrite bool
}

func NewFastFile(parentId, url string) pkg.Upload {
	reg := regexp.MustCompile(`^fast://(\w+):(\d+)/(.+)`)
	params := reg.FindSubmatch([]byte(url))
	size, _ := strconv.ParseInt(string(params[2]), 10, 0)
	return &FastFile{
		parentId: parentId,
		name:     string(params[3]),
		fileMd5:  string(params[1]),
		size:     size,
	}
}

func (f *FastFile) ParentId() string {
	return f.parentId
}
func (f *FastFile) LazyCheck() bool {
	return false
}
func (f *FastFile) Name() string {
	return f.name
}
func (f *FastFile) Overwrite() bool {
	return f.overwrite
}
func (f *FastFile) Size() int64 {
	return f.size
}
func (f *FastFile) SliceNum() int {
	return 0
}
func (f *FastFile) FileMD5() string {
	return f.fileMd5
}
func (f *FastFile) SliceMD5() string {
	return f.fileMd5
}
func (f *FastFile) Part(int64) pkg.UploadPart {
	return nil
}
