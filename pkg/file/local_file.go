package file

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/gowsp/cloud189/pkg"
)

type LocalFile struct {
	Prepare  sync.Once
	Exists   bool
	client   pkg.Client
	parentId string
	uploadId string
	path     string
	info     os.FileInfo
	partName []string
	fileMD5  string
	sliceMD5 string
	writed   int
	sliceNum int
}
type FilePath struct {
	FullPath string
	FileInfo fs.FileInfo
}

func NewLocalFile(parentId string, path *FilePath, client pkg.Client) *LocalFile {
	size := path.FileInfo.Size()
	sliceNum := int(math.Ceil(float64(size) / float64(Slice)))
	return &LocalFile{
		parentId: parentId,
		client:   client,
		info:     path.FileInfo,
		path:     path.FullPath,
		sliceNum: sliceNum,
		partName: make([]string, sliceNum),
	}
}
func (f *LocalFile) Upload() {
	file, err := os.Open(f.path)
	if err != nil {
		log.Fatalf("open %v error %v", f.path, err)
	}
	defer file.Close()
	f.md5()
	num := f.SliceNum()
	for i := 0; i < num; i++ {
		f.writed = i
		part := NewFilePart(file, i, f.partName[i])
		err = f.client.Upload(f, part)
		if err != nil {
			return
		}
		if f.Exists {
			break
		}
	}
}
func (f *LocalFile) SetUploadId(uploadId string) {
	f.uploadId = uploadId
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

func (f *LocalFile) md5() {
	file, err := os.Open(f.path)
	if err != nil {
		log.Fatalf("open %v error %v", f.path, err)
	}
	defer file.Close()

	buf := make([]byte, 32*1024)

	count := f.SliceNum()
	global := md5.New()
	if count == 1 {
		io.CopyBuffer(global, file, buf)
		v := global.Sum(nil)
		f.fileMD5 = hex.EncodeToString(v)
		f.sliceMD5 = f.fileMD5
		f.partName[0] = base64.StdEncoding.EncodeToString(v)
		return
	}

	slices := make([]string, count)
	detail := md5.New()

	for i := 0; i < count; i++ {
		offset := int64(i * int(Slice))
		s := io.NewSectionReader(file, offset, Slice)
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
}

type FilePart struct {
	num  int
	name string
	data *io.SectionReader
}

func NewFilePart(file *os.File, num int, name string) *FilePart {
	data := io.NewSectionReader(file, int64(num*Slice), Slice)
	return &FilePart{data: data, num: num, name: name}
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
