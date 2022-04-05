package pkg

import (
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

	Mkdirs(parentId string, path ...string) (map[string]interface{}, error)

	Copy(taget string, src ...File) error

	Move(taget string, src ...File) error

	Delete(src ...File) error

	Rename(file File, newName string) error

	Download(file File, start int64) (*http.Response, error)

	Uploader
}

type App interface {
	Uploader() Uploader

	Login(name, password string) error

	Sign() error

	Space() (Space, error)

	Stat(path string) (File, error)

	List(file File) ([]File, error)

	ListBy(name string) ([]File, error)

	Mkdir(path string, parents bool) error

	Mkdirs(path ...string) error

	Copy(target string, from ...string) error

	Move(target string, from ...string) error

	Remove(paths ...string) error

	Download(local string, paths ...string) error

	DownloadFile(local string, file File) error

	Upload(cloud string, locals ...string) error
}
