package drive

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/fs"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
)

type LocalFile struct {
	once      sync.Once
	Exists    bool
	client    pkg.Uploader
	parentId  string
	uploadId  string
	path      string
	info      os.FileInfo
	partName  []string
	fileMD5   string
	sliceMD5  string
	writed    int
	sliceNum  int
	overwrite bool
}
type FilePath struct {
	CloudPath string
	LocalPath string
	FileInfo  fs.FileInfo
}

func NewLocalFile(parentId string, path *FilePath, client pkg.Uploader) *LocalFile {
	size := path.FileInfo.Size()
	sliceNum := int(math.Ceil(float64(size) / float64(file.Slice)))
	return &LocalFile{
		parentId: parentId,
		client:   client,
		info:     path.FileInfo,
		path:     path.LocalPath,
		sliceNum: sliceNum,
		partName: make([]string, sliceNum),
	}
}
func (f *LocalFile) Prepare(init func()) {
	f.once.Do(init)
}
func (f *LocalFile) Upload() error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = f.md5()
	if err != nil {
		return nil
	}
	num := f.SliceNum()
	for i := 0; i < num; i++ {
		f.writed = i
		part := NewFilePart(file, i, f.partName[i])
		err = f.client.Upload(f, part)
		if err != nil {
			return err
		}
		if f.Exists {
			break
		}
	}
	return nil
}
func (f *LocalFile) SetUploadId(uploadId string) {
	f.uploadId = uploadId
}
func (f *LocalFile) Overwrite() bool {
	return f.overwrite
}
func (f *LocalFile) SetExists(exists bool) {
	f.Exists = exists
}
func (f *LocalFile) Type() string {
	return "LOCALFILE"
}
func (f *LocalFile) ParentId() string {
	return f.parentId
}
func (f *LocalFile) IsExists() bool {
	return f.Exists
}
func (f *LocalFile) IsComplete() bool {
	return f.writed >= f.sliceNum-1
}
func (f *LocalFile) UploadId() string {
	return f.uploadId
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
	local, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer local.Close()

	buf := make([]byte, 32*1024)

	count := f.SliceNum()
	global := md5.New()
	if count == 1 {
		io.CopyBuffer(global, local, buf)
		v := global.Sum(nil)
		f.fileMD5 = hex.EncodeToString(v)
		f.sliceMD5 = f.fileMD5
		f.partName[0] = base64.StdEncoding.EncodeToString(v)
		return nil
	}

	slices := make([]string, count)
	detail := md5.New()

	for i := 0; i < count; i++ {
		offset := int64(i * int(file.Slice))
		s := io.NewSectionReader(local, offset, file.Slice)
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
		f.partName[i] = base64.StdEncoding.EncodeToString(v)
		detail.Reset()
	}
	slice := strings.Join(slices, "\n")
	detail.Write([]byte(slice))

	f.fileMD5 = hex.EncodeToString(global.Sum(nil))
	f.sliceMD5 = hex.EncodeToString(detail.Sum(nil))
	return nil
}

type FilePart struct {
	num  int
	name string
	data *io.SectionReader
}

func NewFilePart(f *os.File, num int, name string) *FilePart {
	data := io.NewSectionReader(f, int64(num*file.Slice), file.Slice)
	return &FilePart{data: data, num: num, name: name}
}
func (f *FilePart) Name() string {
	return f.name
}
func (f *FilePart) Num() int {
	return int(f.num)

}
func (f *FilePart) Data() io.Reader {
	buf := new(bytes.Buffer)
	io.Copy(buf, f.data)
	return buf
}
