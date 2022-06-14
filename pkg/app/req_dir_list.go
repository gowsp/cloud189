package app

import (
	"github.com/gowsp/cloud189/pkg/file"
)

func (c *api) ListDir(id string) (result []*file.FileInfo, err error) {
	return c.listFile(file.FileType_Dir, id, 1)
}
