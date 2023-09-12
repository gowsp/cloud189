package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/gowsp/cloud189/pkg/util"
)

type Upload struct {
	session *invoker.Session
}

func (client *api) Uploader() pkg.ReadWriter {
	return &Upload{session: client.conf.Session}
}
func (client *Upload) Write(upload pkg.Upload) error {
	data, err := client.init(upload)
	if err != nil {
		return err
	}
	if data.IsExists() {
		return client.commit(upload, data.UploadFileId, "0")
	}
	count := upload.SliceNum()
	for i := 0; i < count; i++ {
		part := upload.Part(int64(i))
		err = client.uploadPart(part, data.UploadFileId)
		if err != nil {
			return err
		}
	}
	return client.commit(upload, data.UploadFileId, "1")
}

func (up *Upload) create(method, u string, f url.Values) *http.Request {
	c := time.Now().Format(time.RFC1123)
	r := util.Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := up.session.Secret

	e := util.EncodeParam(f)
	data := util.AesEncrypt([]byte(e), []byte(l[0:16]))
	h := strings.ToUpper(hex.EncodeToString(data))

	req, err := http.NewRequest(http.MethodGet, "https://upload.cloud.189.cn"+u+"?params="+h, nil)
	if err != nil {
		return nil
	}
	a := make(url.Values)
	a.Set("SessionKey", up.session.Key)
	a.Set("Operate", method)
	a.Set("RequestURI", u)
	a.Set("Date", c)
	a.Set("params", h)

	g := util.Sha1(util.EncodeParam(a), l)
	req.Header.Set("IsJson", "1")
	req.Header.Set("SessionKey", up.session.Key)
	req.Header.Set("Signature", g)
	req.Header.Set("Date", c)
	req.Header.Set("X-Request-ID", r)
	return req
}

func (up *Upload) do(req func() *http.Request, retry int, result interface{}) error {
	r := req()
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(&result)
	}
	if os.Getenv("DEBUG") == "1" {
		rdata, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(rdata))
		data, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(data))
	}
	if retry > 5 {
		return err
	}
	time.Sleep(time.Second)
	return up.do(req, retry+1, result)
}

func (i *Upload) Get(path string, params url.Values, result interface{}) error {
	return i.do(func() *http.Request {
		return i.create(http.MethodGet, path, params)
	}, 0, result)
}
func (i *Upload) Post(path string, params url.Values, result interface{}) error {
	return i.do(func() *http.Request {
		return i.create(http.MethodPost, path, params)
	}, 0, result)
}

type uploadInfo struct {
	UploadType     int    `json:"uploadType,omitempty"`
	UploadHost     string `json:"uploadHost,omitempty"`
	UploadFileId   string `json:"uploadFileId,omitempty"`
	FileDataExists int    `json:"fileDataExists,omitempty"`
}

func (i *uploadInfo) IsExists() bool {
	return i.FileDataExists == 1
}

type initResp struct {
	Code string     `json:"code,omitempty"`
	Data uploadInfo `json:"data,omitempty"`
}

func (r *initResp) GetCode() string {
	return r.Code
}

func (c *Upload) init(i pkg.Upload) (*uploadInfo, error) {
	params := make(url.Values)
	params.Set("parentFolderId", i.ParentId())
	params.Set("fileName", i.Name())
	params.Set("fileSize", strconv.FormatInt(i.Size(), 10))
	params.Set("sliceSize", strconv.Itoa(file.Slice))

	if i.LazyCheck() {
		params.Set("lazyCheck", "1")
	} else {
		params.Set("fileMd5", i.FileMD5())
		params.Set("sliceMd5", i.SliceMD5())
	}
	params.Set("extend", `{"opScene":"1","relativepath":"","rootfolderid":""}`)
	var upload initResp
	if err := c.Get("/person/initMultiUpload", params, &upload); err != nil {
		return nil, err
	}
	if upload.Data.UploadFileId == "" {
		return nil, errors.New("error get upload fileid")
	}
	return &upload.Data, nil
}

type urlResp struct {
	Code string `json:"code,omitempty"`
	Data map[string]struct {
		RequestURL    string `json:"requestURL,omitempty"`
		RequestHeader string `json:"requestHeader,omitempty"`
	} `json:"uploadUrls,omitempty"`
}

func (client *Upload) uploadPart(part pkg.UploadPart, fileId string) error {
	p := make(url.Values)
	num := strconv.Itoa(part.Num() + 1)
	p.Set("partInfo", fmt.Sprintf("%s-%s", num, part.Name()))
	p.Set("uploadFileId", fileId)

	var urlResp urlResp
	if err := client.Get("/person/getMultiUploadUrls", p, &urlResp); err != nil {
		return err
	}
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
	return nil
}

type uploadResult struct {
	Code string `json:"code,omitempty"`
	File struct {
		Id         string `json:"userFileId,omitempty"`
		FileSize   int64  `json:"file_size,omitempty"`
		FileName   string `json:"file_name,omitempty"`
		FileMd5    string `json:"file_md_5,omitempty"`
		CreateDate string `json:"create_date,omitempty"`
	} `json:"file,omitempty"`
}

func (r *uploadResult) GetCode() string {
	return r.Code
}

func (client *Upload) commit(i pkg.Upload, fileId, lazyCheck string) error {
	var result uploadResult
	params := make(url.Values)
	params.Set("fileMd5", i.FileMD5())
	params.Set("sliceMd5", i.SliceMD5())
	params.Set("lazyCheck", lazyCheck)
	params.Set("uploadFileId", fileId)
	if i.Overwrite() {
		params.Set("opertype", "3")
	}
	return client.Get("/person/commitMultiUploadFile", params, &result)
}
