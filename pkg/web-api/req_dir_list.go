package web

import (
	"encoding/json"
	"net/url"
	"os"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

type folder struct {
	FileId   json.Number `json:"id,omitempty"`
	FileName string      `json:"name,omitempty"`
	ParentId string      `json:"pId,omitempty"`
}

func (f *folder) Id() string         { return f.FileId.String() }
func (f *folder) PId() string        { return f.ParentId }
func (f *folder) Name() string       { return f.FileName }
func (f *folder) Size() int64        { return 0 }
func (f *folder) Mode() os.FileMode  { return os.ModePerm }
func (f *folder) ModTime() time.Time { return time.Now() }
func (f *folder) IsDir() bool        { return true }
func (f *folder) Sys() any           { return nil }

func (c *Api) ListDir(id string) ([]pkg.File, error) {
	params := make(url.Values)
	params.Set("id", id)
	params.Set("orderBy", "1")
	params.Set("order", "ASC")

	var folders []*folder
	err := c.invoker.Post("/portal/getObjectFolderNodes.action", params, &folders)
	files := make([]pkg.File, 0)
	for _, f := range folders {
		files = append(files, f)
	}
	return files, err
}
