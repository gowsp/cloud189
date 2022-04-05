package file

import (
	"fmt"

	"github.com/gowsp/cloud189/pkg"
)

type MediaType int

const (
	ALL MediaType = iota
	Pict
	MUSIC
	VIDEO
	DOCUMENT
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
