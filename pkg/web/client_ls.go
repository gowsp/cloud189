package web

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type fileList struct {
	Count   uint32       `json:"count,omitempty"`
	Files   []*CloudFile `json:"fileList,omitempty"`
	Folders []*CloudFile `json:"folderList,omitempty"`
}
type listResp struct {
	Code string   `json:"res_code,omitempty"`
	Data fileList `json:"fileListAo,omitempty"`
}

type searchResp struct {
	ErrorCode string       `json:"errorCode,omitempty"`
	Code      string       `json:"res_code,omitempty"`
	Count     int          `json:"count,omitempty"`
	Files     []*CloudFile `json:"fileList,omitempty"`
	Folders   []*CloudFile `json:"folderList,omitempty"`
}

func (client *Client) Ls(path string) {
	file := client.find(path)
	if file == nil {
		log.Fatalf("%s Not found", path)
	}
	if file.IsFolder {
		files := file.List()
		for _, v := range files {
			if v.IsFolder {
				fmt.Printf("- %s\n", v.Name)
			} else {
				fmt.Printf("f %s\n", v.Name)
			}
		}
	} else {
		fmt.Printf("f %s\n", file.Name)
	}
}

func (client *Client) finds(paths ...string) []*CloudFile {
	files := make([]*CloudFile, 0, len(paths))
	for _, path := range paths {
		file := client.find(path)
		if file == nil {
			log.Printf("%s not found, skip\n", path)
			continue
		}
		files = append(files, file)
	}
	return files
}
func (client *Client) find(cloud string) *CloudFile {
	CheckCloudPath(cloud)
	if cloud == "/" {
		return &Root
	}
	file := &Root
	paths := strings.Split(cloud, "/")
	count := len(paths)
	for i := 1; i < count; i++ {
		path := paths[i]
		if i == 1 {
			if v, f := DefaultDir()[path]; f {
				file = &v
				continue
			}
		}
		files := file.Search(path, false)
		for _, v := range files {
			if v.Name == path {
				file = v
				break
			}
		}
	}
	base := filepath.Base(cloud)
	if base == file.Name {
		return file
	}
	return nil
}

func (client *Client) search(id, name string, includAll bool) []*CloudFile {
	params := make(url.Values)
	params.Set("noCache", fmt.Sprintf("%v", rand.Float64()))
	params.Set("folderId", id)
	params.Set("pageNum", "1")
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
	req.Header.Add("referer", "https://cloud.189.cn/web/main/file/folder/-11")
	var files searchResp
	client.initSesstion()
	resp, err := client.api.Do(req)
	if resp.StatusCode != 200 || err != nil {
		log.Fatalln("list file error")
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&files)
	for _, v := range files.Folders {
		v.IsFolder = true
	}
	if files.Count > 100 {
		log.Println("Too much content, there may be a loss")
	}
	return append(files.Folders, files.Files...)
}

func (client *Client) list(id string) []*CloudFile {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("pageNum", "1")
	params.Set("pageSize", "60")
	params.Set("iconOption", "5")
	params.Set("descending", "true")
	params.Set("orderBy", "lastOpTime")
	params.Set("mediaType", strconv.Itoa(int(ALL)))

	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/open/file/listFiles.action?"+params.Encode(), nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	data, _ := client.api.Do(req)
	var files listResp
	json.NewDecoder(data.Body).Decode(&files)

	for _, v := range files.Data.Folders {
		v.IsFolder = true
	}
	return append(files.Data.Folders, files.Data.Files...)
}
