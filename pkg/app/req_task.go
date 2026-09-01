package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gowsp/cloud189/pkg/invoker"
	"net/url"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/cache"
)

type taskType string

const (
	shareSave taskType = "SHARE_SAVE"
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

type createTaskResponse struct {
	*invoker.NumCodeRsp
	TaskId string `json:"taskId"`
}

func (c *api) createBatchTaskSync(ctx context.Context, taskType taskType, targetFolderId, shareId string, files ...pkg.File) (taskId string, sucCount int, err error) {
	taskId, err = c.createBatchTask(taskType, targetFolderId, shareId, files...)
	if err != nil {
		return
	}

	resp, err := c.waitBatchTask(ctx, taskType, taskId)
	if err != nil {
		return
	}

	return taskId, resp.SuccessedCount, nil
}

func (c *api) createBatchTask(taskType taskType, targetFolderId, shareId string, files ...pkg.File) (taskId string, err error) {
	length := len(files)
	if length == 0 {
		return "", nil
	}
	rm := make([]taskInfo, length)
	for i, v := range files {
		rm[i] = taskInfo{
			Id:    v.Id(),
			Name:  v.Name(),
			IsDir: v.IsDir(),
		}
	}
	data, err := json.Marshal(rm)
	if err != nil {
		return "", err
	}
	params := make(url.Values)
	params.Set("type", string(taskType))
	params.Set("taskInfos", string(data))
	params.Set("targetFolderId", targetFolderId)
	if shareId != "" {
		params.Set("shareId", shareId)
	}

	response := new(createTaskResponse)
	err = c.invoker.Post("/batch/createBatchTask.action", params, response)
	if err != nil {
		return "", err
	}
	if !response.IsSuccess() {
		return "", fmt.Errorf("createBatchTask %s error: %s", taskType, response.Error())
	}

	switch taskType {
	case shareSave:
		cache.Invalid(files...)
		cache.InvalidId(targetFolderId)
	}
	return response.TaskId, nil
}

func (c *api) waitBatchTask(ctx context.Context, taskType taskType, taskId string) (response *CheckTaskResponse, err error) {
	for {
		response, err = c.checkBatchTask(taskType, taskId)
		if err != nil {
			return nil, err
		}
		switch response.TaskStatus {
		case 4:
			return response, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			// 继续下一次循环
		}
	}
}

type CheckTaskResponse struct {
	*invoker.NumCodeRsp
	ErrorCode      string `json:"errorCode"`
	TaskId         string `json:"taskId"`
	TaskStatus     int    `json:"taskStatus"`
	Process        int    `json:"process"`
	FailedCount    int    `json:"failedCount"`
	SkipCount      int    `json:"skipCount"`
	SubTaskCount   int    `json:"subTaskCount"`
	SuccessedCount int    `json:"successedCount"`
}

var taskConflictError = errors.New("there is a conflict with the target object")

func (c *api) checkBatchTask(taskType taskType, taskId string) (response *CheckTaskResponse, err error) {
	response = new(CheckTaskResponse)
	err = c.invoker.Post("/batch/checkBatchTask.action", url.Values{
		"type":   []string{string(taskType)},
		"taskId": []string{taskId},
	}, response)
	if err != nil {
		return response, err
	}
	if !response.IsSuccess() {
		return response, fmt.Errorf("checkBatchTask %s error: %s", taskType, response.Error())
	}

	if response.ErrorCode != "" {
		return response, fmt.Errorf("%s", response.ErrorCode) // eg: InsufficientStorageSpace
	}

	switch response.TaskStatus {
	case 2:
		return response, taskConflictError
	case 3, 4: // 3(进行中), 4(已完成)
		return response, nil
	}

	return response, fmt.Errorf("task failed %d", response.TaskStatus)
}

type getConflictTaskInfoResponse struct {
	*invoker.NumCodeRsp
	SessionKey     string      `json:"sessionKey"`
	TargetFolderId int64       `json:"targetFolderId"`
	TaskId         string      `json:"taskId"`
	TaskInfos      []*TaskInfo `json:"taskInfos"`
	TaskType       int         `json:"taskType"`
}

type TaskInfo struct {
	FileId     int64  `json:"fileId"`
	FileName   string `json:"fileName"`
	IsConflict int    `json:"isConflict"`
	IsFolder   int    `json:"isFolder"`
}

type ManageTaskInfo struct {
	*TaskInfo
	DealWay DealWay `json:"dealWay"`
}

type DealWay int // 1:忽略 2:保留两者 3:替换

const (
	DealWayCancel DealWay = iota
	DealWayIgnore
	DealWayKeepBoth
	DealWayReplace
)

func (ctr *getConflictTaskInfoResponse) manageTaskInfos(dealWay DealWay) []*ManageTaskInfo {
	ret := make([]*ManageTaskInfo, 0, len(ctr.TaskInfos))
	for _, v := range ctr.TaskInfos {
		if v.IsConflict == 1 {
			ret = append(ret, &ManageTaskInfo{
				TaskInfo: v,
				DealWay:  dealWay,
			})
		}
	}

	return ret
}

func (c *api) getConflictTaskInfo(taskType taskType, taskId string) (response *getConflictTaskInfoResponse, err error) {
	response = new(getConflictTaskInfoResponse)
	err = c.invoker.Post("/batch/getConflictTaskInfo.action", url.Values{
		"type":   []string{string(taskType)},
		"taskId": []string{taskId},
	}, response)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		return response, fmt.Errorf("getConflictTaskInfo %s error: %s", taskType, response.Error())
	}

	return response, nil
}

type manageTaskResponse struct {
	*invoker.NumCodeRsp
	Success bool `json:"success"`
}

// 暂只处理冲突文件
func (c *api) manageBatchTask(taskType taskType, taskId, targetFolderId string, manageTaskInfos []*ManageTaskInfo) (err error) {

	data, err := json.Marshal(manageTaskInfos)
	if err != nil {
		return err
	}

	response := new(manageTaskResponse)
	err = c.invoker.Post("/batch/manageBatchTask.action", url.Values{
		"type":           []string{string(taskType)},
		"taskId":         []string{taskId},
		"targetFolderId": []string{targetFolderId},
		"taskInfos":      []string{string(data)},
	}, response)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return fmt.Errorf("manageBatchTask %s error: %s", taskType, response.Error())
	}

	if !response.IsSuccess() {
		return errors.New("manageBatchTask failed")
	}

	return nil
}
