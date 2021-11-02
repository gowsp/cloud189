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
)

type folderInfo struct {
	Path        []fileInfo `json:"path,omitempty"`
	Data        []fileInfo `json:"data,omitempty"`
	PageNum     int        `json:"pageNum,omitempty"`
	PageSize    int        `json:"pageSize,omitempty"`
	RecordCount int        `json:"recordCount,omitempty"`
}

func (f *folderInfo) Name() string {
	i := len(f.Path)
	return f.Path[i-1].Name
}

type fileInfo struct {
	Id          json.Number `json:"fileId,omitempty"`
	ParentId    json.Number `json:"parentId,omitempty"`
	IsFolder    bool        `json:"isFolder,omitempty"`
	Name        string      `json:"fileName,omitempty"`
	Size        int64       `json:"fileSize,omitempty"`
	DownloadUrl string      `json:"downloadUrl,omitempty"`
}

func (f *fileInfo) getDownloadUrl() string {
	return "https:" + f.DownloadUrl
}

func (c *Client) Download(local string, clouds ...string) {
	CheckCloudPath(clouds...)
	for _, cloud := range clouds {
		file := c.find(cloud)
		if file == nil {
			log.Printf("%s not found, skip download\n", cloud)
			continue
		}
		if file.IsFolder {
			c.downByFolderId(file.Id.String(), local, 1)
		} else {
			c.downByFileId(file.Id.String(), local)
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
		err = os.Mkdir(local, 0666)
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
func (client *Client) downByFileId(id, local string) {
	fileInfo, err := client.getFileInfo(id)
	if err != nil {
		log.Println(err)
		return
	}
	client.downFile(fileInfo, local)
}
func (client *Client) downFile(fileInfo *fileInfo, local string) {
	var file *os.File
	stat, err := os.Stat(local)
	if err == nil && stat.IsDir() {
		file, err = os.OpenFile(local+"/"+fileInfo.Name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
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
	if localSize == fileInfo.Size {
		//TODO MD5一致性检查
		return
	}
	req, _ := http.NewRequest(http.MethodGet, fileInfo.getDownloadUrl(), nil)
	// TODO 大文件分片下载
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", localSize, fileInfo.Size))
	resp, err := client.api.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
}

func (client *Client) getFileInfo(id string) (*fileInfo, error) {
	params := make(url.Values)
	params.Set("fileId", id)

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/getFileInfo.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	resp, err := client.api.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var file fileInfo
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
