package web

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/gowsp/cloud189-cli/pkg/file"
)

func (client *Client) Rename(src, dest string) error {
	file.CheckPath(src, dest)
	f, err := client.Stat(src)
	if err != nil {
		return os.ErrNotExist
	}
	dest = path.Base(dest)
	if f.IsDir() {
		client.renameFolder(f.Id(), dest)
	}
	client.renameFile(f.Id(), dest)
	return nil
}

func (client *Client) renameFolder(id, dest string) *file.FileInfo {
	params := make(url.Values)
	params.Set("folderId", id)
	params.Set("destFolderName", dest)
	resp, err := client.api.PostForm("https://cloud.189.cn/api/open/file/renameFolder.action", params)
	if err != nil {
		log.Println()
	}
	defer resp.Body.Close()
	var info file.FileInfo
	json.NewDecoder(resp.Body).Decode(&info)
	return &info
}

func (client *Client) renameFile(id, dest string) *file.FileInfo {
	params := make(url.Values)
	params.Set("fileId", id)
	params.Set("destFileName", dest)
	resp, err := client.api.PostForm("https://cloud.189.cn/api/open/file/renameFile.action", params)
	if err != nil {
		log.Println()
	}
	defer resp.Body.Close()
	var info file.FileInfo
	json.NewDecoder(resp.Body).Decode(&info)
	return &info
}
