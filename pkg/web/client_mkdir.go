package web

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gowsp/cloud189-cli/pkg/file"
)

func (client *Client) Mkdir(clouds ...string) error {
	file.CheckPath(clouds...)
	client.mkdir(file.Root.Id.String(), clouds...)
	return nil
}

func (client *Client) findOrCreateDir(cloud string) folderResp {
	resp := client.mkdir(file.Root.Id.String(), cloud)
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
	var errorResp errorResp
	var result map[string]folderResp
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)
	decoder.Decode(&errorResp)
	defer resp.Body.Close()
	return result
}
