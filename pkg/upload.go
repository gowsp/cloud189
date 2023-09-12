package pkg

import (
	"io"
)

type UploadFile interface {
	Prepare(init func())
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
	Overwrite() bool
	SetExists(exists bool)
	SetUploadId(uploadId string)
}

type UploadPart interface {
	Num() int
	Name() string
	Data() io.Reader
}

type Uploader1 interface {
	Upload(file UploadFile, part UploadPart) error
}
