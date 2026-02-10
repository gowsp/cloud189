package drive

import (
	"context"
	"fmt"

	"github.com/gowsp/cloud189/pkg/file"
)

func (f *FS) SaveShare(shareURL string, target string) error {
	// 1. 解析分享 URL
	shareCode, _, err := ParseShareURL(shareURL)
	if err != nil {
		return fmt.Errorf("parse share URL: %w", err)
	}

	// 2. 获取分享信息
	shareFile, accessCode, shareMode, err := f.api.GetShareInfo(shareCode)
	if err != nil {
		return fmt.Errorf("get share info: %w", err)
	}

	shareId := shareFile.PId()

	// 3. 列出分享内容
	files, err := f.api.ListShareDir(shareFile.Id(), shareId, accessCode, shareMode, shareFile.IsDir(), true)
	if err != nil {
		return fmt.Errorf("list share dir: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("share is empty")
	}

	// 4. 解析目标路径
	dest, err := f.stat(target)
	if err != nil {
		return fmt.Errorf("stat target %s: %w", target, err)
	}
	if !dest.IsDir() {
		return file.ErrFileIsDir
	}

	// 5. 同步保存（DealWay=1 忽略冲突）
	sucCount, err := f.api.ShareSaveSync(context.Background(), 1, dest.Id(), shareId, files...)
	if err != nil {
		return fmt.Errorf("share save: %w", err)
	}

	fmt.Printf("save success, %d files saved to %s\n", sucCount, target)
	return nil
}
