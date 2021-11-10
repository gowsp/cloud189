package file

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
	"math"
	"path"
	"strings"
	"sync"

	"github.com/gowsp/cloud189-cli/pkg"
	"golang.org/x/net/webdav"
)

func NewStreamFile(name string, size int64, client pkg.Client) webdav.File {
	dir := path.Dir(name)
	parent, err := client.Stat(dir)
	if err != nil {
		return nil
	}
	return NewStreamFileWithParent(name, parent.Id(), size, client)
}

func NewStreamFileWithParent(name, parentId string, size int64, client pkg.Client) webdav.File {
	buffSize := Slice
	if size < Slice {
		buffSize = int(size)
	}
	num := int(math.Ceil(float64(size) / float64(Slice)))
	return &StreamFile{
		client:   client,
		partNum:  num,
		parts:    make([]string, num),
		parentId: parentId,
		name:     name,
		size:     size,
		fileMd5:  md5.New(),
		partData: &PartData{
			hash: md5.New(),
			data: bytes.NewBuffer(make([]byte, 0, buffSize)),
		},
	}
}

type StreamFile struct {
	Prepare  sync.Once
	Exists   bool
	client   pkg.Client
	parentId string
	name     string
	size     int64
	partNum  int
	partData *PartData
	fileMd5  hash.Hash
	md5Cache string
	parts    []string
	fileId   string
	writed   int64
}

func (f *StreamFile) Type() string {
	return "STREAM"
}
func (f *StreamFile) ParentId() string {
	return f.parentId
}
func (f *StreamFile) IsExists() bool {
	return f.Exists
}
func (f *StreamFile) IsComplete() bool {
	return f.writed >= f.size
}
func (f *StreamFile) UploadId() string {
	return f.fileId
}
func (f *StreamFile) SetUploadId(fileId string) {
	f.fileId = fileId
}
func (f *StreamFile) Name() string {
	return path.Base(f.name)
}
func (f *StreamFile) Size() int64 {
	return f.size
}
func (f *StreamFile) SliceNum() int {
	return f.partNum
}
func (f *StreamFile) Part() *PartData {
	return f.partData
}
func (f *StreamFile) FileMD5() string {
	if len(f.md5Cache) == 0 {
		v := f.fileMd5.Sum(nil)
		f.md5Cache = hex.EncodeToString(v)
	}
	return f.md5Cache
}
func (f *StreamFile) SliceMD5() string {
	if f.SliceNum() == 1 {
		return f.FileMD5()
	}
	detail := md5.New()
	data := strings.Join(f.parts, "\n")
	detail.Write([]byte(data))
	return hex.EncodeToString(detail.Sum(nil))
}
func (f *StreamFile) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (f *StreamFile) Write(p []byte) (n int, err error) {
	n, _ = f.fileMd5.Write(p)
	offset := f.partData.Writed() + n
	switch {
	case offset > Slice:
		offset = offset - Slice
		f.partData.Write(p[:offset])
		err = f.upload(f.partData)
		if err != nil {
			return offset, err
		}
		f.partData.Reset()
		_, err = f.partData.Write(p[offset:])
	case offset == Slice:
		n, _ = f.partData.Write(p)
		err = f.upload(f.partData)
		if err != nil {
			return n, err
		}
		f.partData.Reset()
	default:
		n, err = f.partData.Write(p)
	}
	f.writed += int64(n)
	if f.IsComplete() {
		err = f.upload(f.partData)
	}
	return
}
func (f *StreamFile) upload(data *PartData) error {
	if f.SliceNum() > 1 {
		f.parts[data.num] = strings.ToUpper(data.MD5())
	}
	return f.client.Upload(f, data)
}
func (f *StreamFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (f *StreamFile) Readdir(count int) ([]fs.FileInfo, error) {
	return make([]fs.FileInfo, 0), nil
}
func (f *StreamFile) Stat() (fs.FileInfo, error) {
	return &FileInfo{MD5: f.FileMD5()}, nil
}
func (f *StreamFile) Close() error {
	return nil
}

type PartData struct {
	data   *bytes.Buffer
	hash   hash.Hash
	md5    []byte
	writed int
	num    int
}

func (b *PartData) Data() io.Reader {
	return b.data
}
func (b *PartData) Write(p []byte) (n int, err error) {
	b.hash.Write(p)
	n, err = b.data.Write(p)
	b.writed += n
	return
}
func (b *PartData) Name() string {
	if len(b.md5) == 0 {
		b.md5 = b.hash.Sum(nil)
	}
	return base64.StdEncoding.EncodeToString(b.md5)
}
func (b *PartData) Num() int {
	return b.num
}
func (b *PartData) MD5() string {
	if len(b.md5) == 0 {
		b.md5 = b.hash.Sum(nil)
	}
	return hex.EncodeToString(b.md5)
}
func (b *PartData) Writed() int {
	return b.writed
}
func (b *PartData) Reset() {
	b.num += 1
	b.md5 = nil
	b.writed = 0
	b.data.Reset()
	b.hash.Reset()
}
