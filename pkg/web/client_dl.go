package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/file"
)

type folderInfo struct {
	Path        []StatInfo `json:"path,omitempty"`
	Data        []StatInfo `json:"data,omitempty"`
	PageNum     int        `json:"pageNum,omitempty"`
	PageSize    int        `json:"pageSize,omitempty"`
	RecordCount int        `json:"recordCount,omitempty"`
}

func (f *folderInfo) Name() string {
	i := len(f.Path)
	return f.Path[i-1].Name()
}

type StatInfo struct {
	Id          json.Number `json:"fileId,omitempty"`
	FileName    string      `json:"fileName,omitempty"`
	FileSize    int64       `json:"fileSize,omitempty"`
	IsFolder    bool        `json:"isFolder,omitempty"`
	FileModTime int64       `json:"lastOpTime,omitempty"`
	DownloadUrl string      `json:"downloadUrl,omitempty"`
}

func (f *StatInfo) Name() string       { return f.FileName }
func (f *StatInfo) Size() int64        { return f.FileSize }
func (f *StatInfo) Mode() os.FileMode  { return 0666 }
func (f *StatInfo) ModTime() time.Time { return time.UnixMilli(f.FileModTime) }
func (f *StatInfo) IsDir() bool        { return f.IsFolder }
func (f *StatInfo) Sys() interface{}   { return nil }
func (f *StatInfo) getDownloadUrl() string {
	return "https:" + f.DownloadUrl
}

func (c *Client) Dl(local string, clouds ...string) {
	file.CheckPath(clouds...)
	for _, cloud := range clouds {
		file, err := c.Stat(cloud)
		if err != nil {
			log.Printf("%s not found, skip download\n", cloud)
			continue
		}
		if file.IsDir() {
			c.downByFolderId(file.Id(), local, 1)
		} else {
			c.Get(file.Id(), local)
		}
	}
}
func (c *Client) downByFolderId(dir string, local string, pageNum int) {
	folder, err := c.getFolderInfo(dir, pageNum)
	if err != nil {
		log.Println(err)
		return
	}
	local = local + "/" + folder.Name()
	_, err = os.Stat(local)
	if os.IsNotExist(err) {
		err = os.Mkdir(local, 0766)
		if err != nil {
			log.Println(err)
			return
		}
	}
	for _, f := range folder.Data {
		if f.IsFolder {
			c.downByFolderId(f.Id.String(), local, 1)
		} else {
			c.downFile(&f, local)
		}
	}
	if folder.RecordCount > folder.PageNum*folder.PageSize {
		c.downByFolderId(dir, local, pageNum+1)
	}
}
func (client *Client) Get(id, local string) error {
	fileInfo, err := client.getFileInfo(id)
	if err != nil {
		return err
	}
	client.downFile(fileInfo, local)
	return nil
}
func (client *Client) downFile(fileInfo *StatInfo, local string) {
	var file *os.File
	stat, err := os.Stat(local)
	if err == nil && stat.IsDir() {
		file, err = os.OpenFile(local+"/"+fileInfo.Name(), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(local, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	}
	if err != nil {
		log.Println(err)
		return
	}
	stat, _ = file.Stat()
	defer file.Close()
	localSize := stat.Size()
	if localSize == fileInfo.Size() {
		//TODO MD5一致性检查
		return
	}
	req, _ := http.NewRequest(http.MethodGet, fileInfo.getDownloadUrl(), nil)
	// TODO 大文件分片下载
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", localSize, fileInfo.Size()))
	resp, err := client.api.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
}
func (client *Client) getFileInfo(id string) (*StatInfo, error) {
	if v, b := file.DefaultIdDir()[id]; b {
		return &StatInfo{Id: v.Id, FileName: v.Name, IsFolder: true}, nil
	}

	params := make(url.Values)
	params.Set("fileId", id)

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/getFileInfo.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	resp, err := client.api.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var file StatInfo
	json.NewDecoder(resp.Body).Decode(&file)
	return &file, nil
}

func (client *Client) getFolderInfo(id string, pageNum int) (*folderInfo, error) {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("pageNum", strconv.Itoa(pageNum))

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/listFiles.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	resp, err := client.api.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var file folderInfo
	json.NewDecoder(resp.Body).Decode(&file)
	return &file, nil
}
