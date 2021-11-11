package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gowsp/cloud189-cli/pkg"
	"github.com/gowsp/cloud189-cli/pkg/file"
)

type fileList struct {
	Count   int              `json:"count,omitempty"`
	Files   []*file.FileInfo `json:"fileList,omitempty"`
	Folders []*file.FileInfo `json:"folderList,omitempty"`
}
type listResp struct {
	Code json.Number `json:"res_code,omitempty"`
	Data fileList    `json:"fileListAo,omitempty"`
}

type searchResp struct {
	ErrorCode string           `json:"errorCode,omitempty"`
	Code      string           `json:"res_code,omitempty"`
	Count     int              `json:"count,omitempty"`
	Files     []*file.FileInfo `json:"fileList,omitempty"`
	Folders   []*file.FileInfo `json:"folderList,omitempty"`
}

func (client *Client) Ls(path string) {
	info, err := client.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}
	if info.IsDir() {
		files := client.list(info.Id(), 1)
		for _, v := range files {
			if v.IsDir() {
				fmt.Printf("- %s\n", v.Name())
			} else {
				fmt.Printf("f %s\n", v.Name())
			}
		}
	} else {
		fmt.Printf("f %s\n", info.Name())
	}
}

func (client *Client) finds(paths ...string) []*file.FileInfo {
	files := make([]*file.FileInfo, 0, len(paths))
	for _, path := range paths {
		info, err := client.Stat(path)
		if err != nil {
			fmt.Printf("%s not found, skip\n", path)
			continue
		}
		files = append(files, info.(*file.FileInfo))
	}
	return files
}
func (client *Client) Stat(cloud string) (info pkg.FileInfo, err error) {
	file.CheckPath(cloud)
	info = &file.FileInfo{FileId: file.Root.Id, FileName: file.Root.Name, IsFolder: true}
	if cloud == "/" {
		return info, nil
	}
	dir := filepath.Dir(cloud)
	if dir == "/" {
		name := filepath.Base(cloud)
		info, err = client.search(info.Id(), name, 1, false)
		return
	}
	paths := strings.Split(dir, "/")
	count := len(paths)
	for i := 1; i < count; i++ {
		path := paths[i]
		info, err = client.findFolder(info.Id(), path)
		if err != nil {
			return nil, err
		}
	}
	if strings.HasSuffix(cloud, "/") {
		return
	}
	name := filepath.Base(cloud)
	info, err = client.search(info.Id(), name, 1, false)
	return
}

func (client *Client) search(id, name string, page int, includAll bool) (pkg.FileInfo, error) {
	params := make(url.Values)
	params.Set("noCache", fmt.Sprintf("%v", rand.Float64()))
	params.Set("folderId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	params.Set("filename", name)
	if includAll {
		params.Set("recursive", "1")
	} else {
		params.Set("recursive", "0")
	}
	params.Set("iconOption", "5")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/open/file/searchFiles.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")
	resp, err := client.api.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var files searchResp
	json.NewDecoder(resp.Body).Decode(&files)
	if files.ErrorCode == "InvalidSessionKey" {
		client.refresh()
		return client.search(id, name, page, includAll)
	}
	for _, folder := range files.Folders {
		folder.IsFolder = true
		if folder.Name() == name {
			return folder, nil
		}
	}
	for _, file := range files.Files {
		if file.Name() == name {
			return file, nil
		}
	}
	if files.Count > len(files.Files)+len(files.Folders) {
		return client.search(id, name, page+1, includAll)
	}
	return nil, fs.ErrNotExist
}

func (Client *Client) findFolder(parentId, name string) (pkg.FileInfo, error) {
	folders, err := client.findFolders(parentId)
	if err != nil {
		return nil, err
	}
	for _, folder := range folders {
		if folder.FileName == name {
			return folder, nil
		}
	}
	return nil, fs.ErrNotExist
}
func (Client *Client) findFolders(parentId string) ([]*file.FileInfo, error) {
	params := make(url.Values)
	params.Set("id", parentId)
	params.Set("orderBy", "1")
	params.Set("order", "ASC")
	resp, err := client.api.PostForm("https://cloud.189.cn/api/portal/getObjectFolderNodes.action", params)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var folders []*file.FileInfo
	json.Unmarshal(data, &folders)
	if len(folders) == 0 && client.isInvalidSession(data) {
		return client.findFolders(parentId)
	}
	return folders, nil
}

func (client *Client) Readdir(id string, count int) []fs.FileInfo {
	data := client.list(id, 1)
	result := make([]fs.FileInfo, len(data))
	for i, v := range data {
		result[i] = v
	}
	return result
}
func (client *Client) list(id string, page int) []*file.FileInfo {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("pageNum", strconv.Itoa(page))
	params.Set("pageSize", "100")
	params.Set("iconOption", "5")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")
	params.Set("mediaType", strconv.Itoa(int(file.ALL)))

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/open/file/listFiles.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	resp, _ := client.api.Do(req)
	body, _ := io.ReadAll(resp.Body)
	var list listResp
	json.Unmarshal(body, &list)
	if list.Code.String() == "" && client.isInvalidSession(body) {
		return client.list(id, page)
	}

	data := list.Data
	for _, v := range data.Folders {
		v.IsFolder = true
	}
	result := make([]*file.FileInfo, 0, data.Count)
	result = append(result, data.Files...)
	result = append(result, data.Folders...)
	if data.Count > len(data.Files)+len(data.Folders) {
		result = append(result, client.list(id, page+1)...)
	}
	return result
}
