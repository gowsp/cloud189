package web

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189-cli/pkg"
	"github.com/gowsp/cloud189-cli/pkg/file"
	"github.com/gowsp/cloud189-cli/pkg/util"
)

type uploadResp interface {
	GetCode() string
}
type uploadInfo struct {
	UploadType     int    `json:"uploadType,omitempty"`
	UploadHost     string `json:"uploadHost,omitempty"`
	UploadFileId   string `json:"uploadFileId,omitempty"`
	FileDataExists int    `json:"fileDataExists,omitempty"`
}
type initResp struct {
	Code string     `json:"code,omitempty"`
	Data uploadInfo `json:"data,omitempty"`
}

func (r *initResp) GetCode() string {
	return r.Code
}

type urlResp struct {
	Code string                `json:"code,omitempty"`
	Data map[string]uploadUrls `json:"uploadUrls,omitempty"`
}

func (r *urlResp) GetCode() string {
	return r.Code
}

type uploadUrls struct {
	RequestURL    string `json:"requestURL,omitempty"`
	RequestHeader string `json:"requestHeader,omitempty"`
}

func (client *Client) Up(cloud string, locals ...string) {
	file.CheckPath(cloud)
	dir := client.findOrCreateDir(cloud)
	for _, local := range locals {
		i := file.NewLocalFile(dir.Id.String(), local, client)
		i.Upload()
	}
}
func (client *Client) init(i pkg.UploadFile, parentId string) *uploadInfo {
	f := make(url.Values)
	f.Set("parentFolderId", parentId)
	f.Set("fileName", i.Name())
	f.Set("fileSize", strconv.FormatInt(i.Size(), 10))
	f.Set("sliceSize", strconv.Itoa(file.Slice))

	if i.SliceNum() > 1 {
		f.Set("lazyCheck", "1")
	} else {
		f.Set("fileMd5", i.FileMD5())
		f.Set("sliceMd5", i.SliceMD5())
	}
	var upload initResp
	client.uploadRequest(func() *http.Request {
		return client.createRequest("/person/initMultiUpload", f)
	}, &upload)
	return &upload.Data
}

func (client *Client) check(i pkg.UploadFile, fileId string) *uploadInfo {
	var upload initResp
	params := make(url.Values)
	params.Set("fileMd5", i.FileMD5())
	params.Set("sliceMd5", i.SliceMD5())
	params.Set("uploadFileId", fileId)
	client.uploadRequest(func() *http.Request {
		return client.createRequest("/person/checkTransSecond", params)
	}, &upload)
	return &upload.Data
}

func (client *Client) Upload(upload pkg.UploadFile, part pkg.UploadPart) error {
	switch upload.Type() {
	case "STREAM":
		file := upload.(*file.StreamFile)
		file.Prepare.Do(func() {
			info := client.init(file, file.ParentId())

			file.SetUploadId(info.UploadFileId)
			if file.UploadId() == "" {
				log.Fatalln("error get upload fileid")
			}
			file.Exists = info.FileDataExists == 1
		})
	case "LOCALFILE":
		file := upload.(*file.LocalFile)
		file.Prepare.Do(func() {
			info := client.init(file, file.ParentId())

			fileId := info.UploadFileId
			file.SetUploadId(fileId)
			if fileId == "" {
				log.Fatalln("error get upload fileid")
			}
			if info.FileDataExists == 1 {
				file.Exists = true
				return
			}
			info = client.check(file, fileId)
			if info.FileDataExists == 1 {
				file.Exists = true
				return
			}
		})
	}
	if upload.IsExists() {
		log.Println("file exists, fast upload")
		client.commit(upload, upload.UploadId(), "0")
		return nil
	}
	err := client.UploadPart(part, upload.UploadId())
	if err != nil {
		return err
	}
	if upload.IsComplete() {
		client.commit(upload, upload.UploadId(), "1")
	}
	return nil
}
func (client *Client) UploadPart(part pkg.UploadPart, fileId string) error {
	p := make(url.Values)
	num := strconv.Itoa(part.Num() + 1)
	p.Set("partInfo", fmt.Sprintf("%s-%s", num, part.Name()))
	p.Set("uploadFileId", fileId)

	var urlResp urlResp
	client.uploadRequest(func() *http.Request {
		return client.createRequest("/person/getMultiUploadUrls", p)
	}, &urlResp)
	log.Printf("start uploading part %s\n", num)

	upload := urlResp.Data["partNumber_"+num]
	req, _ := http.NewRequest(http.MethodPut, upload.RequestURL, part.Data())
	headers := strings.Split(upload.RequestHeader, "&")
	for _, v := range headers {
		i := strings.Index(v, "=")
		req.Header.Set(v[0:i], v[i+1:])
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload error, %s", string(data))
	}
	log.Printf("part %s upload completed\n", num)
	return nil
}

type UploadResult struct {
	Code string       `json:"code,omitempty"`
	File UploadDetail `json:"file,omitempty"`
}
type UploadDetail struct {
	Id         string `json:"userFileId,omitempty"`
	FileSize   int64  `json:"file_size,omitempty"`
	FileName   string `json:"file_name,omitempty"`
	FileMd5    string `json:"file_md_5,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

func (r *UploadResult) GetCode() string {
	return r.Code
}
func (client *Client) commit(i pkg.UploadFile, fileId, lazyCheck string) UploadDetail {
	var result UploadResult
	params := make(url.Values)
	params.Set("fileMd5", i.FileMD5())
	params.Set("sliceMd5", i.SliceMD5())
	params.Set("lazyCheck", lazyCheck)
	params.Set("uploadFileId", fileId)
	client.uploadRequest(func() *http.Request {
		return client.createRequest("/person/commitMultiUploadFile", params)
	}, &result)
	return result.File
}
func (client *Client) uploadRequest(req func() *http.Request, data uploadResp) {
	resp, err := client.api.Do(req())
	if err != nil {
		log.Fatalln(err)
	}
	json.NewDecoder(resp.Body).Decode(data)
	resp.Body.Close()

	switch data.GetCode() {
	case "SUCCESS":
		return
	case "InvalidSessionKey":
		client.refresh()
		client.uploadRequest(req, data)
	case "InvalidSignature":
		client.uploadRequest(req, data)
	default:
		log.Fatalln(data.GetCode())
	}
}

func (client *Client) createRequest(u string, f url.Values) *http.Request {
	c := strconv.FormatInt(time.Now().UnixMilli(), 10)
	r := util.Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := util.Random("xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx")
	l = l[0 : 16+int(16*rand.Float32())]

	e := util.EncodeParam(f)
	data := util.AesEncrypt([]byte(e), []byte(l[0:16]))
	h := hex.EncodeToString(data)

	sessionKey := client.sesstionKey()
	a := make(url.Values)
	a.Set("SessionKey", sessionKey)
	a.Set("Operate", http.MethodGet)
	a.Set("RequestURI", u)
	a.Set("Date", c)
	a.Set("params", h)
	g := util.SHA1(util.EncodeParam(a), l)

	rsa := client.rsa()
	b := rsa.encrypt(l)

	req, err := http.NewRequest(http.MethodGet, "https://upload.cloud.189.cn"+u+"?params="+h, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("accept", "application/json;charset=UTF-8")
	req.Header.Set("SessionKey", sessionKey)
	req.Header.Set("Signature", hex.EncodeToString(g))
	req.Header.Set("X-Request-Date", c)
	req.Header.Set("X-Request-ID", r)
	req.Header.Set("EncryptionText", base64.StdEncoding.EncodeToString(b))
	req.Header.Set("PkId", rsa.PkId)
	return req
}
