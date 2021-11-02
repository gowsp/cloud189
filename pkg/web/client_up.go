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
	"os"
	"strconv"
	"strings"
	"time"

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

const (
	KB = 1 << 10
	MB = 1 << 20

	Slice = 10 * MB
)

func (client *Client) Upload(cloud string, locals ...string) {
	CheckCloudPath(cloud)
	client.initSesstion()
	dir := client.findOrCreateDir(cloud)
	for _, local := range locals {
		client.upload(dir.Id.String(), local)
	}
}
func (client *Client) upload(parentId, name string) {
	i := NewUploadFile(name, Slice)
	info := client.init(i, parentId)

	fileId := info.UploadFileId
	if fileId == "" {
		log.Fatalln("error get upload fileid")
	}
	if info.FileDataExists == 1 {
		log.Println("file exists, fast upload")
		client.commit(i, fileId, "0")
		return
	}

	var upload initResp
	params := make(url.Values)
	params.Set("uploadFileId", fileId)
	client.sendRequest(func() *http.Request {
		return client.createRequest("/person/getUploadedPartsInfo", params)
	}, &upload)

	info = client.check(i, fileId)
	if info.FileDataExists == 1 {
		log.Println("file exists, fast upload")
		client.commit(i, fileId, "1")
		return
	}
	client.copy(i, fileId)
	client.commit(i, fileId, "1")
}

func (client *Client) init(i *UploadFile, parentId string) *uploadInfo {
	f := make(url.Values)
	f.Set("parentFolderId", parentId)
	f.Set("fileName", i.BaseName())
	f.Set("fileSize", strconv.FormatInt(i.Size(), 10))
	f.Set("sliceSize", strconv.Itoa(Slice))

	a := i.SliceNum()
	if a > 1 {
		f.Set("lazyCheck", "1")
	} else {
		f.Set("fileMd5", i.FileMD5())
		f.Set("sliceMd5", i.SliceMD5())
	}
	var upload initResp
	client.sendRequest(func() *http.Request {
		return client.createRequest("/person/initMultiUpload", f)
	}, &upload)
	return &upload.Data
}

func (client *Client) check(i *UploadFile, fileId string) *uploadInfo {
	var upload initResp
	params := make(url.Values)
	params.Set("fileMd5", i.FileMD5())
	params.Set("sliceMd5", i.SliceMD5())
	params.Set("uploadFileId", fileId)
	client.sendRequest(func() *http.Request {
		return client.createRequest("/person/checkTransSecond", params)
	}, &upload)
	return &upload.Data
}

func (client *Client) copy(u *UploadFile, fileId string) {
	name := u.FullName()
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("open %v error %v", name, err)
	}
	defer file.Close()

	for i, part := range u.Parts() {
		p := make(url.Values)
		num := strconv.Itoa(i + 1)
		p.Set("partInfo", fmt.Sprintf("%s-%s", num, part.Name))
		p.Set("uploadFileId", fileId)

		var urlResp urlResp
		client.sendRequest(func() *http.Request {
			return client.createRequest("/person/getMultiUploadUrls", p)
		}, &urlResp)
		log.Printf("start uploading part %s\n", num)

		upload := urlResp.Data["partNumber_"+num]
		s := io.NewSectionReader(file, part.Offset, u.SliceSize)
		req, _ := http.NewRequest(http.MethodPut, upload.RequestURL, s)
		headers := strings.Split(upload.RequestHeader, "&")
		for _, v := range headers {
			i := strings.Index(v, "=")
			req.Header.Set(v[0:i], v[i+1:])
		}
		resp, _ := http.DefaultClient.Do(req)
		if resp.StatusCode != 200 {
			data, _ := io.ReadAll(resp.Body)
			log.Fatalf("upload error, %s\n", string(data))
		}
		resp.Body.Close()
		log.Printf("part %s upload completed\n", num)
	}
	log.Println("file upload completed")
}
func (client *Client) commit(i *UploadFile, fileId, lazyCheck string) {
	var upload initResp
	params := make(url.Values)
	params.Set("fileMd5", i.FileMD5())
	params.Set("sliceMd5", i.SliceMD5())
	params.Set("lazyCheck", lazyCheck)
	params.Set("uploadFileId", fileId)
	client.sendRequest(func() *http.Request {
		return client.createRequest("/person/commitMultiUploadFile", params)
	}, &upload)
}
func (client *Client) sendRequest(req func() *http.Request, data uploadResp) {
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
		client.sendRequest(req, data)
	case "InvalidSignature":
		client.sendRequest(req, data)
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
