package drive

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ParseShareURL 从天翼云盘分享链接中解析出 shareCode 和 accessCode。
// 支持的格式：
//   - https://cloud.189.cn/t/reiQruJJnyEv
//   - https://cloud.189.cn/t/reiQruJJnyEv（访问码：v5oy）
//   - reiQruJJnyEv
//   - reiQruJJnyEv（访问码：v5oy）
func ParseShareURL(raw string) (shareCode, accessCode string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", fmt.Errorf("empty share URL")
	}

	// 提取访问码（中文括号格式）
	accessCodeRe := regexp.MustCompile(`（访问码：([^）]+)）`)
	if matches := accessCodeRe.FindStringSubmatch(raw); len(matches) == 2 {
		accessCode = matches[1]
		raw = strings.TrimSpace(accessCodeRe.ReplaceAllString(raw, ""))
	}

	// 如果是完整 URL，提取路径最后一段
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, parseErr := url.Parse(raw)
		if parseErr != nil {
			return "", "", fmt.Errorf("invalid share URL: %w", parseErr)
		}
		parts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
		if len(parts) == 0 {
			return "", "", fmt.Errorf("cannot extract share code from URL: %s", raw)
		}
		shareCode = parts[len(parts)-1]
	} else {
		shareCode = raw
	}

	if shareCode == "" {
		return "", "", fmt.Errorf("cannot extract share code from: %s", raw)
	}

	// 如果有访问码，拼接为 GetShareInfo 期望的格式
	if accessCode != "" {
		shareCode = shareCode + "（访问码：" + accessCode + "）"
	}

	return shareCode, accessCode, nil
}
