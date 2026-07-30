package cmd

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

type JSONFileEntry struct {
	Path         string `json:"path"`
	RelativePath string `json:"relative_path"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	FileID       string `json:"file_id"`
	ParentFileID string `json:"parent_file_id"`
	ModifiedAt   string `json:"modified_at"`
	Checksum     string `json:"checksum"`
}

type JSONSpace struct {
	OK        bool   `json:"ok"`
	Capacity  uint64 `json:"capacity"`
	Available uint64 `json:"available"`
}

type JSONError struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func fileToJSONEntry(cloudPath string, info os.FileInfo) JSONFileEntry {
	return fileToJSONEntryWithBase(cloudPath, path.Dir(cleanCloudPath(cloudPath)), info)
}

func fileToJSONEntryWithBase(cloudPath, basePath string, info os.FileInfo) JSONFileEntry {
	cloudPath = cleanCloudPath(cloudPath)
	modifiedAt := ""
	if !info.ModTime().IsZero() {
		modifiedAt = info.ModTime().UTC().Format(time.RFC3339)
	}
	entry := JSONFileEntry{
		Path:         cloudPath,
		RelativePath: relativeCloudPath(basePath, cloudPath),
		Name:         info.Name(),
		Size:         info.Size(),
		IsDir:        info.IsDir(),
		ModifiedAt:   modifiedAt,
		Checksum:     checksumFromFileInfo(info),
	}
	if file, ok := info.(interface {
		Id() string
		PId() string
	}); ok {
		entry.FileID = file.Id()
		entry.ParentFileID = file.PId()
	}
	return entry
}

func writeJSON(v any) error {
	return writeJSONTo(os.Stdout, v)
}

func writeJSONError(err error) error {
	if err == nil {
		return nil
	}
	return writeJSONTo(os.Stderr, JSONError{OK: false, Error: err.Error()})
}

func writeJSONTo(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func cleanCloudPath(name string) string {
	if name == "" {
		return "/"
	}
	if !path.IsAbs(name) {
		name = "/" + name
	}
	return path.Clean(name)
}

func joinCloudPath(parent, name string) string {
	return path.Join(cleanCloudPath(parent), name)
}

func relativeCloudPath(basePath, cloudPath string) string {
	basePath = cleanCloudPath(basePath)
	cloudPath = cleanCloudPath(cloudPath)
	if basePath == cloudPath {
		return path.Base(cloudPath)
	}
	if basePath == "/" {
		return strings.TrimPrefix(cloudPath, "/")
	}
	prefix := strings.TrimSuffix(basePath, "/") + "/"
	if strings.HasPrefix(cloudPath, prefix) {
		return strings.TrimPrefix(cloudPath, prefix)
	}
	return strings.TrimPrefix(cloudPath, "/")
}

func checksumFromFileInfo(info fs.FileInfo) string {
	value := reflect.Indirect(reflect.ValueOf(info))
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return ""
	}
	for _, name := range []string{"Md5", "MD5", "FileMd5"} {
		field := value.FieldByName(name)
		if field.IsValid() && field.Kind() == reflect.String {
			return field.String()
		}
	}
	return ""
}
