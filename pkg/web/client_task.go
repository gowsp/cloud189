package web

import (
	"encoding/json"
	"net/http"

	"github.com/gowsp/cloud189-cli/pkg/file"
)

type taskType string

const (
	COPY   taskType = "COPY"
	MOVE   taskType = "MOVE"
	DELETE taskType = "DELETE"
)

func (client *Client) runTask(taskType taskType, paths ...string) {
	file.CheckPath(paths...)
	if taskType == DELETE {
		files := client.finds(paths...)
		client.createTask(taskType, "", files...)
		return
	}
	length := len(paths)
	dest := paths[length-1]
	src := client.finds(paths[:length-1]...)
	target := client.findOrCreateDir(dest)
	client.createTask(taskType, target.Id.String(), src...)
}

type taskInfo struct {
	Id       string `json:"fileId,omitempty"`
	Name     string `json:"fileName,omitempty"`
	IsFolder uint   `json:"isFolder,omitempty"`
}

func (client *Client) createTask(taskType taskType, targetFolderId string, files ...*file.FileInfo) {
	length := len(files)
	if length == 0 {
		return
	}
	req, _ := http.NewRequest(http.MethodPost, "https://cloud.189.cn/api/open/batch/createBatchTask.action", nil)
	req.Header.Add("accept", "application/json;charset=UTF-8")

	rm := make([]taskInfo, length)
	for i, v := range files {
		isFolder := 0
		if v.IsFolder {
			isFolder = 1
		}
		rm[i] = taskInfo{Id: string(v.FileId), Name: v.Name(), IsFolder: uint(isFolder)}
	}
	data, _ := json.Marshal(rm)

	params := req.URL.Query()
	params.Set("taskInfos", string(data))
	params.Set("type", string(taskType))
	params.Set("targetFolderId", targetFolderId)
	req.URL.RawQuery = params.Encode()
	client.api.Do(req)

}
