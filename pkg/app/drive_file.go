package app

import (
	"encoding/json"
	"io/fs"
	"strings"
	"time"
)

type Time time.Time

func (j *Time) UnmarshalJSON(b []byte) error {
	json := string(b)
	s := strings.Trim(json, "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = Time(t)
	return nil
}

type folder struct {
	ID           json.Number `json:"id"`
	ParentID     json.Number `json:"parentId"`
	FileCata     int         `json:"fileCata"`
	FileCount    int         `json:"fileCount"`
	FileListSize int         `json:"fileListSize"`
	LastOpTime   Time        `json:"lastOpTime"`
	CreateDate   Time        `json:"createDate"`
	DirName      string      `json:"name"`
	Rev          string      `json:"rev"`
	StarLabel    int         `json:"starLabel"`
}

func (f *folder) Info() (fs.FileInfo, error) { return f, nil }

func (f *folder) Id() string         { return f.ID.String() }
func (f *folder) PId() string        { return f.ParentID.String() }
func (f *folder) Name() string       { return f.DirName }
func (f *folder) Size() int64        { return 0 }
func (f *folder) Type() fs.FileMode  { return fs.ModeDir }
func (f *folder) Mode() fs.FileMode  { return fs.ModeDir }
func (f *folder) ModTime() time.Time { return time.Time(f.LastOpTime) }
func (f *folder) IsDir() bool        { return true }
func (f *folder) Sys() any           { return nil }

type fileInfo struct {
	ParentID string
	ID       json.Number `json:"id"`

	Md5         string `json:"md5"`
	MediaType   int    `json:"mediaType"`
	FileCata    int    `json:"fileCata"`
	FileName    string `json:"name"`
	FileSize    int64  `json:"size"`
	Orientation int    `json:"orientation"`
	Rev         string `json:"rev"`
	StarLabel   int    `json:"starLabel"`
	LastOpTime  Time   `json:"lastOpTime"`
	CreateDate  Time   `json:"createDate"`
}

func (f *fileInfo) Info() (fs.FileInfo, error) { return f, nil }

func (f *fileInfo) Id() string         { return f.ID.String() }
func (f *fileInfo) PId() string        { return f.ParentID }
func (f *fileInfo) Name() string       { return f.FileName }
func (f *fileInfo) Size() int64        { return f.FileSize }
func (f *fileInfo) Mode() fs.FileMode  { return fs.ModeType }
func (f *fileInfo) Type() fs.FileMode  { return fs.ModeType }
func (f *fileInfo) ModTime() time.Time { return time.Time(f.LastOpTime) }
func (f *fileInfo) IsDir() bool        { return false }
func (f *fileInfo) Sys() any           { return nil }
