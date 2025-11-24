package app

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/invoker"
)

func (c *api) Detail(id string) (string, error) {
	var info map[string]string
	err := c.invoker.Get("/getFileDownloadUrl.action", url.Values{"fileId": {id}}, &info)
	return info["fileDownloadUrl"], err
}

func (c *api) Download(file pkg.File, start int64) (*http.Response, error) {
	if file.IsDir() {
		return nil, errors.New("not support download dir")
	}
	url, _ := c.Detail(file.Id())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, file.Size()))
	return c.invoker.Send(req)
}

type SimpleFolder struct {
	Fid   int    `json:"fid"`
	Fname string `json:"fname"`
}
type FileDetail struct {
	*invoker.NumCodeRsp
	CreateDate      string `json:"createDate"`
	FileDownloadURL string `json:"fileDownloadUrl"`
	FilePath        string `json:"filePath"`
	Icon            struct {
		LargeURL  string `json:"largeUrl"`
		MediumURL string `json:"mediumUrl"`
		SmallURL  string `json:"smallUrl"`
	} `json:"icon"`
	ID                 int64  `json:"id"`
	LastOpTime         int64  `json:"lastOpTime"`
	LastOpTimeStr      string `json:"lastOpTimeStr"`
	Md5                string `json:"md5"`
	MediaType          int    `json:"mediaType"`
	FileName           string `json:"name"`
	ParentFolderListAO struct {
		ParentFolderList []SimpleFolder `json:"parentFolderList"`
	} `json:"parentFolderListAO"`
	ParentID int64 `json:"parentId"`
	Rev      int64 `json:"rev"`
	FileSize int64 `json:"size"`
}

func (f *FileDetail) Id() string         { return fmt.Sprintf("%d", f.ID) }
func (f *FileDetail) PId() string        { return fmt.Sprintf("%d", f.ParentID) }
func (f *FileDetail) Name() string       { return f.FileName }
func (f *FileDetail) Size() int64        { return f.FileSize }
func (f *FileDetail) Mode() os.FileMode  { return os.ModeSymlink }
func (f *FileDetail) ModTime() time.Time { return time.Unix(f.LastOpTime, 0) }
func (f *FileDetail) IsDir() bool        { return f.Md5 == "" }
func (f *FileDetail) Sys() any           { return f.ParentFolderListAO.ParentFolderList }

type FolderInfo struct {
	CreateDate         string `json:"createDate"`
	CreateTime         int64  `json:"createTime"`
	FileID             int64  `json:"fileId"`
	FileName           string `json:"fileName"`
	FilePath           string `json:"filePath"`
	LastOpTime         int64  `json:"lastOpTime"`
	LastOpTimeStr      string `json:"lastOpTimeStr"`
	ParentFolderListAO struct {
		ParentFolderList []SimpleFolder `json:"parentFolderList"`
	} `json:"parentFolderListAO"`
	ParentID int64 `json:"parentId"`
	Rev      int64 `json:"rev"`
}

type FolderExtInfo struct {
	*invoker.NumCodeRsp
	FileCountNum   uint64 `json:"fileCount"`
	FileSizeNum    uint64 `json:"fileSize"`
	FolderCountNum uint64 `json:"folderCount"`
	FolderID       int64  `json:"folderId"`
	RecursionFlag  int    `json:"recursionFlag"`
	TaskID         string `json:"taskId"`
	TaskStatus     int    `json:"taskStatus"`
}

func (f *FolderExtInfo) FileCount() uint64   { return f.FileCountNum }
func (f *FolderExtInfo) FileSize() uint64    { return f.FileSizeNum }
func (f *FolderExtInfo) FolderCount() uint64 { return f.FolderCountNum }

func (c *api) DirUsage(file pkg.File) (pkg.Usage, error) {
	response := &struct {
		ResCode    int    `json:"res_code"`
		ResMessage string `json:"res_message"`
		TaskId     string `json:"taskId"`
	}{}
	err := c.invoker.Get("/file/createFolderExtInfoTask.action", url.Values{"folderId": {file.Id()}}, &response)
	if err != nil {
		return nil, err
	}
	rsp := new(FolderExtInfo)
	req := url.Values{"taskId": {response.TaskId}}
	// 循环查询任务结果，直到状态变为4（完成）或出现错误
	for {
		err = c.invoker.Get("/file/queryTaskResult.action", req, rsp)
		if err != nil {
			return nil, err
		}
		// 如果任务完成则返回结果
		if rsp.TaskStatus == 4 {
			return rsp, nil
		}
		// 如果不是状态3（进行中），则返回错误
		if rsp.TaskStatus != 3 {
			return nil, fmt.Errorf("unexpected task status: %d", rsp.TaskStatus)
		}
		// 等待1.5秒再重试
		time.Sleep(1500 * time.Millisecond)
	}
}
func (c *api) GetFolderInfoById(id string) (*FolderInfo, error) {
	response := &struct {
		ResCode    int    `json:"res_code"`
		ResMessage string `json:"res_message"`
		*FolderInfo
	}{}
	err := c.invoker.Get("/getFolderInfo.action", url.Values{
		"folderId":   {id},
		"folderPath": {},
		"pathList":   {"1"},
		"dt":         {"3"},
	}, &response)
	return response.FolderInfo, err
}
func (c *api) GetFolderInfoByPath(path string) (*FolderInfo, error) {
	response := &struct {
		ResCode    int    `json:"res_code"`
		ResMessage string `json:"res_message"`
		*FolderInfo
	}{}
	err := c.invoker.Get("/getFolderInfo.action", url.Values{
		"folderId":   {},
		"folderPath": {path},
		"pathList":   {"1"},
		"dt":         {"3"},
	}, &response)
	return response.FolderInfo, err
}

func (c *api) Stat(path string) (pkg.File, error) {
	response := new(FileDetail)
	err := c.invoker.Get("/getFileInfo.action", url.Values{
		"fileId":     {},
		"filePath":   {path},
		"pathList":   {"1"},
		"iconOption": {"0"},
	}, response)
	if rsp, ok := err.(invoker.BadRsp); ok {
		if rsp.IsError(invoker.ErrFileNotFound) {
			return nil, os.ErrNotExist
		}
	}
	return response, err
}
func (c *api) GetFileInfoById(id string) (*FileDetail, error) {
	response := &struct {
		ResCode    int    `json:"res_code"`
		ResMessage string `json:"res_message"`
		*FileDetail
	}{}
	err := c.invoker.Get("/getFileInfo.action", url.Values{
		"fileId":     {id},
		"filePath":   {},
		"pathList":   {"1"},
		"iconOption": {"0"},
	}, &response)
	return response.FileDetail, err
}
