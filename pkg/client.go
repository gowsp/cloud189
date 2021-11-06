package pkg

import (
	"io/fs"
)

type FileInfo interface {
	Id() string
	fs.FileInfo
}

type Client interface {
	Cmd

	Uploader

	Get(id, local string) error

	Mkdir(name ...string) error

	Readdir(id string, count int) []fs.FileInfo

	Stat(name string) (FileInfo, error)

	Rename(oldName, newName string) error

	RemoveAll(name ...string) error
}
