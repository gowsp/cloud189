package file

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

func IsSystem(parent, name string) bool {
	if parent == Root.Id() {
		if _, ok := system[name]; ok {
			return true
		}
	}
	return false
}

func IsSystemDir(file pkg.File) bool {
	if file.Id() == Root.Id() {
		return true
	}
	if file.PId() == Root.Id() {
		if _, ok := system[file.Name()]; ok {
			return true
		}
	}
	return false
}

var Root = &sysFolder{FileId: "-11", FileName: "全部文件"}
var system map[string]string = map[string]string{
	"同步盘":  "0",
	"私密空间": "-10",
	"我的图片": "-12",
	"我的视频": "-13",
	"我的音乐": "-14",
	"我的文档": "-15",
	"我的应用": "-16",
}

type sysFolder struct {
	FileId   json.Number
	FileName string
	ParentId string
}

func (f *sysFolder) Id() string         { return f.FileId.String() }
func (f *sysFolder) PId() string        { return f.ParentId }
func (f *sysFolder) Name() string       { return f.FileName }
func (f *sysFolder) Size() int64        { return 0 }
func (f *sysFolder) Mode() os.FileMode  { return os.ModePerm }
func (f *sysFolder) ModTime() time.Time { return time.Now() }
func (f *sysFolder) IsDir() bool        { return true }
func (f *sysFolder) Sys() any           { return nil }
