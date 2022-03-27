package web

import (
	"encoding/json"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
)

type taskType string

const (
	COPY   taskType = "COPY"
	MOVE   taskType = "MOVE"
	DELETE taskType = "DELETE"
)

type taskInfo struct {
	Id       string `json:"fileId,omitempty"`
	Name     string `json:"fileName,omitempty"`
	IsFolder uint   `json:"isFolder,omitempty"`
}

func (c *Api) Copy(target string, file ...pkg.File) error {
	return c.createTask(COPY, target, file...)
}
func (c *Api) Move(target string, file ...pkg.File) error {
	return c.createTask(MOVE, target, file...)
}
func (c *Api) Delete(file ...pkg.File) error {
	return c.createTask(DELETE, "", file...)
}

func (c *Api) createTask(taskType taskType, targetFolderId string, files ...pkg.File) error {
	length := len(files)
	if length == 0 {
		return nil
	}
	rm := make([]taskInfo, length)
	for i, v := range files {
		isFolder := 0
		if v.IsDir() {
			isFolder = 1
		}
		rm[i] = taskInfo{Id: string(v.Id()), Name: v.Name(), IsFolder: uint(isFolder)}
	}
	data, err := json.Marshal(rm)
	if err != nil {
		return err
	}
	params := make(url.Values)
	params.Set("type", string(taskType))
	params.Set("taskInfos", string(data))
	params.Set("targetFolderId", targetFolderId)
	var f FileInfo
	return c.invoker.Post("/open/batch/createBatchTask.action", params, &f)
}
