package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
)



const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
	TB = 1 << 40

	Slice = 10 * MB
)

func ReadableSize(size uint64) string {
	result := float64(size)
	unit := ""
	switch {
	case size >= TB:
		unit = "T"
		result /= TB
	case size >= GB:
		unit = "G"
		result /= GB
	case size >= MB:
		unit = "M"
		result /= MB
	case size >= KB:
		unit = "K"
		result /= KB
	}
	return fmt.Sprintf("%.2f%s", result, unit)
}

func ReadableFileInfo(info pkg.File) string {
	var size string
	if info.IsDir() {
		size = "-"
	} else {
		size = ReadableSize(uint64(info.Size()))
	}
	modTime := info.ModTime().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%-10s%-22s%s", size, modTime, info.Name())
}

type ModTime time.Time

func (j *ModTime) UnmarshalJSON(b []byte) error {
	json := string(b)
	s := strings.Trim(json, "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = ModTime(t)
	return nil
}

type FileInfo struct {
	IsFolder    bool
	ParentId    json.Number
	FileId      json.Number `json:"id,omitempty"`
	FileName    string      `json:"name,omitempty"`
	FileSize    int64       `json:"size,omitempty"`
	MD5         string      `json:"md5,omitempty"`
	FileModTime ModTime     `json:"lastOpTime,omitempty"`
}

func (f *FileInfo) Id() string        { return f.FileId.String() }
func (f *FileInfo) PId() string       { return f.ParentId.String() }
func (f *FileInfo) Name() string      { return f.FileName }
func (f *FileInfo) Size() int64       { return f.FileSize }
func (f *FileInfo) Mode() os.FileMode { return os.ModePerm }
func (f *FileInfo) ModTime() time.Time {
	return time.Time(f.FileModTime)
}
func (f *FileInfo) IsDir() bool      { return f.IsFolder }
func (f *FileInfo) Sys() interface{} { return nil }
func (f *FileInfo) ContentType(ctx context.Context) (string, error) {
	return path.Ext(f.Name()), nil
}
func (f *FileInfo) ETag(ctx context.Context) (string, error) {
	return strconv.FormatInt(f.ModTime().Unix(), 10), nil
}
