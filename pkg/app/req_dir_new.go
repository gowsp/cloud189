package app

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"

	"github.com/gowsp/cloud189/pkg"
)

type makeDirResp struct {
	ResCode    int    `json:"res_code,omitempty"`
	ResMessage string `json:"res_message,omitempty"`

	Folder *folder
}

func (r *makeDirResp) Error() error {
	if r.ResCode == 0 {
		return nil
	}
	return errors.New(r.ResMessage)
}

func (r *makeDirResp) UnmarshalJSON(b []byte) error {
	resp := &struct {
		ResCode    int    `json:"res_code,omitempty"`
		ResMessage string `json:"res_message,omitempty"`
	}{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	r.ResCode = resp.ResCode
	r.ResMessage = resp.ResMessage
	if resp.ResCode != 0 {
		return nil
	}
	var folder folder
	if err := json.Unmarshal(b, &folder); err != nil {
		return err
	}
	r.Folder = &folder
	return nil
}

func (c *api) Mkdir(parent pkg.File, name string) (pkg.File, error) {
	var result makeDirResp
	dir, base := path.Split(name)
	params := url.Values{"folderName": {base}, "relativePath": {dir}, "parentFolderId": {parent.Id()}}
	err := c.invoker.Post("/createFolder.action", params, &result)
	if err != nil {
		return nil, err
	}
	if err = result.Error(); err != nil {
		return nil, err
	}
	return result.Folder, nil
}
