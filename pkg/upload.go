package pkg

import "io"

type UploadFile interface {
	ParentId() string
	Name() string
	Size() int64
	SliceNum() int
	FileMD5() string
	SliceMD5() string
	IsExists() bool
	Type() string
	IsComplete() bool
	UploadId() string
}

type UploadPart interface {
	Name() string
	Num() int
	Data() io.Reader
}

type Uploader interface {
	Upload(file UploadFile, part UploadPart) error
}
