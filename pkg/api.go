package pkg

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Space interface {
	Available() uint64
	Capacity() uint64
}

type File interface {
	Id() string
	PId() string
	os.FileInfo
}
type FileExt struct {
	FileCount   int64
	CreateTime  time.Time
	DownloadUrl string
}

var Root = &SysFolder{FileId: "-11", FileName: "全部文件"}
var System map[string]string = map[string]string{
	"同步盘":  "0",
	"私密空间": "-10",
	"我的图片": "-12",
	"我的视频": "-13",
	"我的音乐": "-14",
	"我的文档": "-15",
	"我的应用": "-16",
}

type SysFolder struct {
	FileId   json.Number
	FileName string
	ParentId string
}

func (f *SysFolder) Id() string         { return f.FileId.String() }
func (f *SysFolder) PId() string        { return f.ParentId }
func (f *SysFolder) Name() string       { return f.FileName }
func (f *SysFolder) Size() int64        { return 0 }
func (f *SysFolder) Mode() os.FileMode  { return os.ModePerm }
func (f *SysFolder) ModTime() time.Time { return time.Now() }
func (f *SysFolder) IsDir() bool        { return true }
func (f *SysFolder) Sys() any           { return nil }

type Api interface {
	Sign() error

	Space() (Space, error)

	Login(name, password string) error

	Find(id, name string) (File, error)

	FindDir(id, name string) (File, error)

	FindFile(id, name string) (File, error)

	Detail(id string) (File, error)

	ListFile(id string) ([]File, error)

	Mkdir(parentId, path string, parents bool) error

	Mkdirs(parentId string, path ...string) error

	Copy(taget string, src ...File) error

	Move(taget string, src ...File) error

	Delete(src ...File) error

	Rename(file File, newName string) error

	Download(file File, start int64) (*http.Response, error)

	Uploader
}

type App interface {
	Login(name, password string) error

	Sign() error

	Space() (Space, error)

	Stat(path string) (File, error)

	List(file File) ([]File, error)

	Mkdir(path string, parents bool) error

	Mkdirs(path ...string) error

	Copy(target string, from ...string) error

	Move(target string, from ...string) error

	Remove(paths ...string) error

	Download(local string, paths ...string) error

	Upload(cloud string, locals ...string) error
}
