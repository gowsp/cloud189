package drive

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/gowsp/cloud189/pkg"
)

func IsFastFile(name string) bool {
	r := regexp.MustCompile(`^fast://(\w+):(\d+)/(.+)`)
	return r.Match([]byte(name))
}

type FastFile struct {
	once      sync.Once
	client    pkg.Uploader
	parentId  string
	uploadId  string
	name      string
	fileMd5   string
	size      int64
	overwrite bool
}

func NewFastFile(parentId, url string, client pkg.Uploader) *FastFile {
	reg := regexp.MustCompile(`^fast://(\w+):(\d+)/(.+)`)
	params := reg.FindSubmatch([]byte(url))
	size, err := strconv.ParseInt(string(params[2]), 10, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &FastFile{
		parentId: parentId,
		client:   client,
		name:     string(params[3]),
		fileMd5:  string(params[1]),
		size:     size,
	}
}
func (f *FastFile) Upload() {
	f.client.Upload(f, nil)
}
func (f *FastFile) Prepare(init func()) {
	f.once.Do(init)
}
func (f *FastFile) ParentId() string {
	return f.parentId
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
	return 1
}
func (f *FastFile) FileMD5() string {
	return f.fileMd5
}
func (f *FastFile) SliceMD5() string {
	return f.fileMd5
}
func (f *FastFile) SetExists(exists bool) {
}
func (f *FastFile) IsExists() bool {
	return true
}
func (f *FastFile) Type() string {
	return "FAST"
}
func (f *FastFile) IsComplete() bool {
	return true
}
func (f *FastFile) UploadId() string {
	return f.uploadId
}
func (f *FastFile) SetUploadId(uploadId string) {
	f.uploadId = uploadId
}
