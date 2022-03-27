package drive

import (
	"fmt"
	"os"
)

func (client *Client) Upload(cloud string, locals ...string) error {
	dir, err := client.Stat(cloud)
	if len(locals) > 1 && os.IsNotExist(err) {
		client.Mkdir(cloud, true)
		dir, err = client.Stat(cloud)
	}
	for _, local := range locals {
		if IsNetFile(local) {
			f := NewNetFile(dir.Id(), local, client.api)
			f.Upload()
			continue
		}
		if IsFastFile(local) {
			f := NewFastFile(dir.Id(), local, client.api)
			f.Upload()
			continue
		}
		info, err := os.Stat(local)
		if err != nil {
			fmt.Printf("open %v error %v\n", local, err)
			continue
		}
		// if info.IsDir() {
		// 	client.uploadFolder(dir.Id(), local)
		// 	continue
		// }
		l := &FilePath{FullPath: local, FileInfo: info}
		i := NewLocalFile(dir.Id(), l, client.api)
		i.Upload()
	}
	return nil
}

// func (client *Client) uploadFolder(parentId, local string) {
// 	folder := file.ReadDir(local)
// 	client.api.Mkdirs(parentId, folder.Folders...)
// 	for k, v := range resp {
// 		files := folder.Files[k]
// 		if files == nil {
// 			continue
// 		}
// 		for val := files.Front(); val != nil; val = val.Next() {
// 			f := val.Value.(*file.FilePath)
// 			i := NewLocalFile(v.Id(), f, client.api)
// 			i.Upload()
// 		}
// 	}
// 	if folder.Dirict == nil {
// 		return
// 	}
// 	for val := folder.Dirict.Front(); val != nil; val = val.Next() {
// 		f := val.Value.(*FilePath)
// 		i := NewLocalFile(parentId, f, client.api)
// 		i.Upload()
// 	}
// }
