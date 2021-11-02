package web

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
)

type MediaType int

const (
	ALL MediaType = iota
	Pict
	MUSIC
	VIDEO
	DOCUMENT
)

var Sec CloudFile = CloudFile{Id: "-10", Name: "私密空间", IsFolder: true}
var Root CloudFile = CloudFile{Id: "-11", Name: "全部文件", IsFolder: true}
var Sync CloudFile = CloudFile{ParentId: "-11", Id: "0", Name: "同步盘", IsFolder: true}
var Picture CloudFile = CloudFile{ParentId: "-11", Id: "-12", Name: "我的图片", IsFolder: true}
var Vedio CloudFile = CloudFile{ParentId: "-11", Id: "-13", Name: "我的视频", IsFolder: true}
var Music CloudFile = CloudFile{ParentId: "-11", Id: "-14", Name: "我的音乐", IsFolder: true}
var Document CloudFile = CloudFile{ParentId: "-11", Id: "-15", Name: "我的文档", IsFolder: true}
var App CloudFile = CloudFile{ParentId: "-11", Id: "-16", Name: "我的应用", IsFolder: true}

var defaultDir map[string]CloudFile
var defaultDirInstance sync.Once

func DefaultDir() map[string]CloudFile {
	defaultDirInstance.Do(func() {
		defaultDir = map[string]CloudFile{}
		defaultDir["同步盘"] = Sync
		defaultDir["私密空间"] = Sec
		defaultDir["全部文件"] = Root
		defaultDir["我的图片"] = Picture
		defaultDir["我的视频"] = Vedio
		defaultDir["我的音乐"] = Music
		defaultDir["我的文档"] = Document
		defaultDir["我的应用"] = App
	})
	return defaultDir
}

func CheckCloudPath(paths ...string) {
	for _, v := range paths {
		if !strings.HasPrefix(v, "/") {
			log.Fatalf("path %s must start with /\n", v)
		}
	}
}

type CloudFile struct {
	Id       json.Number `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	ParentId json.Number `json:"parentId,omitempty"`
	IsFolder bool        `json:"isFolder,omitempty"`
}

func (f *CloudFile) List() []*CloudFile {
	if f.IsFolder {
		return GetClient().list(string(f.Id))
	}
	return make([]*CloudFile, 0)
}
func (f *CloudFile) Search(name string, includAll bool) []*CloudFile {
	if f.IsFolder {
		return GetClient().search(string(f.Id), name, includAll)
	}
	return make([]*CloudFile, 0)
}
