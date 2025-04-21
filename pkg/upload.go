package pkg

import (
	"errors"
	"io"
	"strings"

	"github.com/gowsp/cloud189/pkg/util"
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

type UploadConfig struct {
	Num    uint32
	Parten string
}

func (c *UploadConfig) NewTask() *util.TaskPool {
	return util.NewTask(int(c.Num))
}
func (c *UploadConfig) Check() (err error) {
	if c.Num <= 0 {
		return errors.New("error number of parallels")
	}
	c.Parten = strings.TrimSpace(c.Parten)
	return nil
}
