package pkg

import (
	"io/fs"
	"net/http"
)

type Drive interface {
	fs.StatFS
	fs.ReadDirFS
	Usage(name string) (Usage, error)
	Space() (Space, error)
	Mkdir(name string) error
	Delete(name ...string) error
	QrLogin() error
	Login(username, password string) error
	Copy(target string, source ...string) error
	Move(target string, source ...string) error
	Upload(config UploadConfig, cloud string, locals ...string) error
	UploadFrom(file Upload) error
	Download(local string, cloud ...string) error
	Share(prifix, cloud string) (func(http.ResponseWriter, *http.Request), error)
	GetDownloadUrl(cloud string) (string, error)
}

type FileType uint16

const (
	ALL FileType = iota
	FILE
	DIR
)

type ReadWriter interface {
	// upload
	Write(info Upload) error
}

type Upload interface {
	ParentId() string
	Name() string
	Size() int64
	SliceNum() int
	FileMD5() string
	SliceMD5() string
	Overwrite() bool
	Part(int64) UploadPart
	LazyCheck() bool
}

type DriveApi interface {
	QrLogin() error

	PwdLogin(username, password string) error

	// get upload
	Uploader() ReadWriter

	// get download link
	Download(file File, start int64) (*http.Response, error)

	// sign for space
	Sign() error

	// get space info
	Space() (Space, error)

	// get file info
	Stat(path string) (File, error)

	// get folder usage
	DirUsage(file File) (Usage, error)

	// searce file by type
	Search(parent File, fileType FileType, name string) ([]File, error)

	// list file by type
	List(parent File, fileType FileType) ([]File, error)

	// mkdir
	Mkdir(parent File, name string) (File, error)

	// rename file
	Rename(target File, name string) error

	// move file
	Move(target File, source ...File) error

	// copy file
	Copy(target File, source ...File) error

	// delete file
	Delete(file ...File) error
}
