package web

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

type taskType string

const (
	COPY   taskType = "COPY"
	MOVE   taskType = "MOVE"
	DELETE taskType = "DELETE"
)

type taskInfo struct {
	Id    string
	Name  string
	IsDir bool
}

func (t *taskInfo) MarshalJSON() ([]byte, error) {
	dir := 0
	if t.IsDir {
		dir = 1
	}
	json := fmt.Sprintf(`{"fileId":"%s","fileName":"%s","isFolder":%d}`, t.Id, t.Name, dir)
	return []byte(json), nil
}

func (c *Api) Copy(target string, files ...pkg.File) error {
	return c.createTask(COPY, target, files...)
}
func (c *Api) Move(target string, files ...pkg.File) error {
	return c.createTask(MOVE, target, files...)
}
func (c *Api) Delete(files ...pkg.File) error {
	return c.createTask(DELETE, "", files...)
}

func (c *Api) createTask(taskType taskType, targetFolderId string, files ...pkg.File) error {
	length := len(files)
	if length == 0 {
		return nil
	}
	rm := make([]taskInfo, length)
	for i, v := range files {
		rm[i] = taskInfo{Id: string(v.Id()), Name: v.Name(), IsDir: v.IsDir()}
	}
	data, err := json.Marshal(rm)
	if err != nil {
		return err
	}
	params := make(url.Values)
	params.Set("type", string(taskType))
	params.Set("taskInfos", string(data))
	params.Set("targetFolderId", targetFolderId)
	var result map[string]interface{}
	err = c.invoker.Post("/open/batch/createBatchTask.action", params, &result)
	if err != nil {
		return err
	}
	if taskType == DELETE {
		cache.Delete(files...)
	}
	return nil
}
