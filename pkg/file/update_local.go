package file

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math"
	"os"
	"strings"

	"github.com/gowsp/cloud189/pkg"
)

type LocalFile struct {
	parentId  string
	file      *os.File
	info      os.FileInfo
	partName  []string
	fileMD5   string
	sliceMD5  string
	sliceNum  int
	partNum   int
	overwrite bool
}

func NewLocalFile(parentId string, path string) pkg.Upload {
	source, err := os.Open(path)
	if err != nil {
		return nil
	}
	info, err := source.Stat()
	if err != nil {
		return nil
	}
	size := info.Size()
	sliceNum := int(math.Ceil(float64(size) / float64(Slice)))
	return &LocalFile{
		parentId: parentId,
		info:     info,
		file:     source,
		sliceNum: sliceNum,
		partName: make([]string, sliceNum),
	}
}
func (f *LocalFile) Close() {
	f.file.Close()
}
func (f *LocalFile) Part(num int64) pkg.UploadPart {
	f.partNum = int(num)
	data := io.NewSectionReader(f.file, int64(num*Slice), Slice)
	buff := bytes.NewBuffer(nil)
	io.Copy(buff, data)
	return &FilePart{data: buff, num: num, name: f.partName[num]}
}

func (f *LocalFile) Overwrite() bool {
	return f.overwrite
}
func (f *LocalFile) ParentId() string {
	return f.parentId
}
func (f *LocalFile) Name() string {
	return f.info.Name()
}
func (f *LocalFile) Size() int64 {
	return f.info.Size()
}
func (f *LocalFile) SliceNum() int {
	return f.sliceNum
}
func (f *LocalFile) LazyCheck() bool {
	return false
}
func (f *LocalFile) FileMD5() string {
	if f.fileMD5 == "" {
		f.md5()
	}
	return f.fileMD5
}
func (f *LocalFile) SliceMD5() string {
	if f.sliceMD5 == "" {
		f.md5()
	}
	return f.sliceMD5
}

func (f *LocalFile) md5() error {
	count := f.sliceNum
	slices := make([]string, count)

	global := md5.New()
	detail := md5.New()
	writer := io.MultiWriter(global, detail)

	buf := make([]byte, 32*1024)
	for i := 0; i < count; i++ {
		offset := int64(i * int(Slice))
		s := io.NewSectionReader(f.file, offset, Slice)
		io.CopyBuffer(writer, s, buf)
		v := detail.Sum(nil)
		slices[i] = strings.ToUpper(hex.EncodeToString(v))
		f.partName[i] = base64.StdEncoding.EncodeToString(v)
		detail.Reset()
	}
	f.fileMD5 = hex.EncodeToString(global.Sum(nil))
	if count > 1 {
		slice := strings.Join(slices, "\n")
		detail.Write([]byte(slice))
		f.sliceMD5 = hex.EncodeToString(detail.Sum(nil))
	} else {
		f.sliceMD5 = f.fileMD5
	}
	return nil
}

type FilePart struct {
	num  int64
	name string
	data *bytes.Buffer
}

func (f *FilePart) Name() string {
	return f.name
}
func (f *FilePart) Num() int {
	return int(f.num)
}
func (f *FilePart) Data() io.Reader {
	return f.data
}
