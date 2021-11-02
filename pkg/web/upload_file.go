package web

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

type UploadFile struct {
	name      string
	info      os.FileInfo
	parts     []UploadPart
	fileMD5   string
	sliceMD5  string
	SliceSize int64
}
type UploadPart struct {
	Name   string
	Offset int64
}

func NewUploadFile(name string, sliceSize int64) *UploadFile {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatalf("open %v error %v", name, err)
	}
	if info.IsDir() {
		log.Fatalf("folder %s uploads are not supported", name)
	}
	return &UploadFile{name: name, info: info, SliceSize: sliceSize}
}
func (f *UploadFile) FullName() string {
	return f.name
}
func (f *UploadFile) BaseName() string {
	return f.info.Name()
}
func (f *UploadFile) Size() int64 {
	return f.info.Size()
}
func (f *UploadFile) SliceNum() int {
	return int(math.Ceil(float64(f.info.Size()) / float64(f.SliceSize)))
}
func (f *UploadFile) FileMD5() string {
	if f.fileMD5 == "" {
		f.md5()
	}
	return f.fileMD5
}
func (f *UploadFile) SliceMD5() string {
	if f.sliceMD5 == "" {
		f.md5()
	}
	return f.sliceMD5
}
func (f *UploadFile) Parts() []UploadPart {
	if len(f.parts) == 0 {
		f.md5()
	}
	return f.parts
}
func (f *UploadFile) md5() {
	file, err := os.Open(f.name)
	if err != nil {
		log.Fatalf("open %v error %v", f.name, err)
	}
	defer file.Close()

	count := f.SliceNum()

	buf := make([]byte, 32*1024)

	global := md5.New()
	f.parts = make([]UploadPart, count)
	if count == 1 {
		io.CopyBuffer(global, file, buf)
		v := global.Sum(nil)
		n := base64.StdEncoding.EncodeToString(v)

		f.fileMD5 = hex.EncodeToString(v)
		f.sliceMD5 = f.fileMD5
		f.parts[0] = UploadPart{Name: n}
		return
	}
	slices := make([]string, count)
	detail := md5.New()

	for i := 0; i < count; i++ {
		offset := int64(i * int(f.SliceSize))
		s := io.NewSectionReader(file, offset, f.SliceSize)
		r := bufio.NewReader(s)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				detail.Write(buf[0:n])
				global.Write(buf[0:n])
			}
			if err != nil {
				break
			}
		}
		v := detail.Sum(nil)
		slices[i] = strings.ToUpper(hex.EncodeToString(v))
		n := base64.StdEncoding.EncodeToString(v)
		f.parts[i] = UploadPart{Name: n, Offset: offset}
		detail.Reset()
	}
	slice := strings.Join(slices, "\n")
	detail.Write([]byte(slice))

	f.fileMD5 = hex.EncodeToString(global.Sum(nil))
	f.sliceMD5 = hex.EncodeToString(detail.Sum(nil))
}
