package app

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/invoker"
)

type shareInfo struct {
	*invoker.NumCodeRsp
	AccessCode     string  `json:"accessCode"`
	Creator        Creator `json:"creator"`
	ExpireTime     int     `json:"expireTime"`
	ExpireType     int     `json:"expireType"`
	FileCreateDate string  `json:"fileCreateDate"`
	FileId         string  `json:"fileId"`
	FileLastOpTime Time    `json:"fileLastOpTime"`
	FileName       string  `json:"fileName"`
	FileSize       int64   `json:"fileSize"`
	FileType       string  `json:"fileType"`
	IsFolder       bool    `json:"isFolder"`
	NeedAccessCode int     `json:"needAccessCode"`
	ReviewStatus   int     `json:"reviewStatus"`
	ShareDate      int64   `json:"shareDate"`
	ShareId        int64   `json:"shareId"`
	ShareMode      int     `json:"shareMode"`
	ShareType      int     `json:"shareType"`
}

type Creator struct {
	IconURL      string `json:"iconURL"`
	NickName     string `json:"nickName"`
	Oper         bool   `json:"oper"`
	OwnerAccount string `json:"ownerAccount"`
	SuperVip     int    `json:"superVip"`
	Vip          int    `json:"vip"`
}

func (f *shareInfo) Info() (fs.FileInfo, error) { return f, nil }

func (f *shareInfo) Id() string   { return f.FileId }
func (f *shareInfo) PId() string  { return fmt.Sprintf("%d", f.ShareId) }
func (f *shareInfo) Name() string { return f.FileName }
func (f *shareInfo) Size() int64  { return f.FileSize }
func (f *shareInfo) Mode() fs.FileMode {
	if f.IsFolder {
		return fs.ModeDir
	} else {
		return fs.ModeType
	}
}
func (f *shareInfo) Type() fs.FileMode  { return f.Mode() }
func (f *shareInfo) ModTime() time.Time { return time.Time(f.FileLastOpTime) }
func (f *shareInfo) IsDir() bool {
	if f.IsFolder {
		return true
	} else {
		return false
	}
}
func (f *shareInfo) Sys() any { return nil }

func (c *api) GetShareInfo(shareCode string) (result pkg.File, accessCode string, shareMode int, err error) {
	ret, err := c.getShareInfo(shareCode)
	if err != nil {
		return nil, "", 0, err
	}
	return ret, ret.AccessCode, ret.ShareMode, nil
}

func (c *api) getShareInfo(shareCode string) (result *shareInfo, err error) {
	response := new(shareInfo)
	err = c.invoker.Get("/open/share/getShareInfoByCodeV2.action", url.Values{
		"shareCode": {shareCode},
	}, response)
	if rsp, ok := err.(invoker.BadRsp); ok {
		if rsp.IsError(invoker.ErrFileNotFound) {
			return nil, os.ErrNotExist
		}
	}
	if !response.IsSuccess() {
		return nil, fmt.Errorf("getShareInfo error: %s", response.Error())
	}

	return response, err
}

type ShareDirInfo struct {
	*invoker.NumCodeRsp
	ExpireTime int            `json:"expireTime"`
	ExpireType int            `json:"expireType"`
	List       *ShareFileList `json:"fileListAO"`
	LastRev    int            `json:"lastRev"`
}

type ShareFileList struct {
	Count        int            `json:"count"`
	FileList     []*shareFile   `json:"fileList"`
	FileListSize int            `json:"fileListSize"`
	FolderList   []*shareFolder `json:"folderList"`
}
type shareFile struct {
	fileInfo
	Icon struct {
		LargeUrl string `json:"largeUrl"`
		SmallUrl string `json:"smallUrl"`
	} `json:"icon"`
}

type shareFolder struct {
	folder
	*ShareFileList //Recursive
}

func (list *ShareFileList) Files() []pkg.File {
	if list == nil {
		return []pkg.File{}
	}
	files := make([]pkg.File, 0, len(list.FileList)+len(list.FolderList))
	for _, folder := range list.FolderList {
		files = append(files, folder)
	}
	for _, file := range list.FileList {
		files = append(files, file)
	}
	return files
}
func (list *ShareFileList) Tree() string {
	if list == nil {
		return ""
	}

	var sb strings.Builder
	list.writeTree("", &sb)
	return sb.String()
}

func (list *ShareFileList) writeTree(prefix string, sb *strings.Builder) {
	folders := list.FolderList
	files := list.FileList

	// Sort both lists individually
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].DirName < folders[j].DirName
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].FileName < files[j].FileName
	})

	folderIdx, fileIdx := 0, 0
	total := len(folders) + len(files)

	for i := 0; i < total; i++ {
		var isFolder bool
		var name string
		var nextFolder *shareFolder

		// Decide whether to pick from folders or files
		if folderIdx < len(folders) && fileIdx < len(files) {
			if folders[folderIdx].DirName < files[fileIdx].FileName {
				isFolder = true
			} else {
				isFolder = false
			}
		} else if folderIdx < len(folders) {
			isFolder = true
		} else {
			isFolder = false
		}

		if isFolder {
			nextFolder = folders[folderIdx]
			name = nextFolder.DirName
			folderIdx++
		} else {
			name = files[fileIdx].FileName
			fileIdx++
		}

		isLast := i == total-1
		connector := "├── "
		childPrefix := "│   "
		if isLast {
			connector = "└── "
			childPrefix = "    "
		}

		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(name)
		sb.WriteString("\n")

		if isFolder && nextFolder.FileList != nil {
			nextFolder.ShareFileList.writeTree(prefix+childPrefix, sb)
		}
	}
}

type ShareTaskOption struct {
	Context context.Context
	Sync    bool
	DealWay DealWay
}

var (
	defaultShareTaskOption = &ShareTaskOption{
		Context: context.Background(),
		Sync:    true,
		DealWay: DealWayIgnore,
	}
)

func (c *api) ListShareDir(fileId, shareId, accessCode string, shareMode int, isFolder, recursive bool) ([]pkg.File, error) {
	result, err := c.listShareDirWithRecursive(fileId, shareId, accessCode, shareMode, isFolder, recursive)
	if err != nil {
		return nil, err
	}
	return result.Files(), nil
}

func (c *api) listShareDir(fileId, shareId, accessCode string, shareMode int, isFolder bool) (*ShareDirInfo, error) {
	response := new(ShareDirInfo)
	err := c.invoker.Get("/open/share/listShareDir.action", url.Values{
		"pageNum":        {"1"},
		"pageSize":       {"200"},
		"fileId":         {fileId},
		"shareDirFileId": {fileId},
		"isFolder":       {fmt.Sprintf("%v", isFolder)},
		"shareId":        {shareId},
		"shareMode":      {fmt.Sprintf("%d", shareMode)},
		"iconOption":     {"5"},
		"orderBy":        {"lastOpTime"},
		"descending":     {"true"},
		"accessCode":     {accessCode},
	}, response)
	if rsp, ok := err.(invoker.BadRsp); ok {
		if rsp.IsError(invoker.ErrFileNotFound) {
			return nil, os.ErrNotExist
		}
	}
	if !response.IsSuccess() {
		return nil, fmt.Errorf("listShareDir error: %s", response.Error())
	}

	return response, err
}

func (c *api) listShareDirWithRecursive(fileId, shareId, accessCode string, shareMode int, isFolder, recursive bool) (*ShareFileList, error) {
	sdi, err := c.listShareDir(fileId, shareId, accessCode, shareMode, isFolder)
	if err != nil {
		return nil, err
	}

	if !recursive {
		return sdi.List, err
	}

	if len(sdi.List.FolderList) > 0 {
		for _, folder := range sdi.List.FolderList {
			folder.ShareFileList, err = c.listShareDirWithRecursive(folder.ID.String(), shareId, accessCode, shareMode, isFolder, recursive)
			if err != nil {
				return nil, fmt.Errorf("listShareDirRecursively \"%s\" error: %w", folder.DirName, err)
			}
		}
	}

	return sdi.List, nil
}

func (c *api) ShareSave(targetFolderId, shareId string, files ...pkg.File) (taskId string, err error) {
	return c.createBatchTask(shareSave, targetFolderId, shareId, files...)
}
func (c *api) ShareSaveSync(ctx context.Context, dealWay int, targetFolderId, shareId string, files ...pkg.File) (sucCount int, err error) {
	if DealWay(dealWay) == DealWayCancel {
		_, sucCount, err = c.createBatchTaskSync(context.Background(), shareSave, targetFolderId, shareId, files...)
	} else {
		_, sucCount, err = c.shareSaveSyncWithDealWay(ctx, DealWay(dealWay), targetFolderId, shareId, files...)
	}
	return
}

// dealWay(1:忽略 2:保留两者 3:替换)
func (c *api) shareSaveSyncWithDealWay(ctx context.Context, dealWay DealWay, targetFolderId, shareId string, files ...pkg.File) (taskId string, sucCount int, err error) {
	taskId, sucCount, err = c.createBatchTaskSync(ctx, shareSave, targetFolderId, shareId, files...)
	if err == nil {
		return
	}

	// 冲突时获取冲突信息并修改任务
	if err == taskConflictError {
		resp, err := c.getConflictTaskInfo(shareSave, taskId)
		if err != nil {
			return taskId, sucCount, err
		}

		err = c.manageBatchTask(shareSave, taskId, targetFolderId, resp.manageTaskInfos(dealWay))
		if err != nil {
			return taskId, sucCount, err
		}

		taskResp, err := c.waitBatchTask(ctx, shareSave, taskId)
		if err != nil {
			return taskId, sucCount, err
		}
		return taskId, taskResp.SuccessedCount, nil
	}
	return
}
