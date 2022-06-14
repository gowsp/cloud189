package file

import "strconv"

type MediaType int

const (
	ALL MediaType = iota
	Pict
	MUSIC
	VIDEO
	DOCUMENT
)

type FileType int

func (f FileType) String() string {
	return strconv.Itoa(int(f))
}

const (
	FileType_All FileType = iota
	FileType_File
	FileType_Dir
)
