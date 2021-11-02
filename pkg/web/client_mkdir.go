package web

import (
	"encoding/json"
	"log"
	"net/url"
)

func (client *Client) Mkdir(clouds ...string) {
	CheckCloudPath(clouds...)
	client.initSesstion()
	client.mkdir(Root.Id.String(), clouds...)
}

func (client *Client) findOrCreateDir(cloud string) folderResp {
	resp := client.mkdir(Root.Id.String(), cloud)
	target := resp[cloud[1:]]
	if target.Success {
		return target
	}
	log.Fatalf("find or create dir %s error", cloud)
	return folderResp{}
}

type folderList struct {
	ParentId string   `json:"parentId,omitempty"`
	Paths    []string `json:"paths,omitempty"`
}
type folderResp struct {
	Id      json.Number `json:"result,omitempty"`
	Success bool        `json:"success,omitempty"`
}

func (client *Client) mkdir(parentId string, paths ...string) map[string]folderResp {
	for i, v := range paths {
		paths[i] = v[1:]
	}
	f := folderList{ParentId: parentId, Paths: paths}
	data, _ := json.Marshal(f)
	params := make(url.Values)
	params.Set("folderList", string(data))
	resp, err := client.api.PostForm("https://cloud.189.cn/v2/createFolders.action", params)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]folderResp
	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()
	return result
}
