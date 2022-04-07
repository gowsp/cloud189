package app

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/invoker"
)

type userSession struct {
	LoginName           string `json:"loginName,omitempty"`
	SessionKey          string `json:"sessionKey,omitempty"`
	SessionSecret       string `json:"sessionSecret,omitempty"`
	KeepAlive           int    `json:"keepAlive,omitempty"`
	FileDiffSpan        int    `json:"getFileDiffSpan,omitempty"`
	UserInfoSpan        int    `json:"getUserInfoSpan,omitempty"`
	FamilySessionKey    string `json:"familySessionKey,omitempty"`
	FamilySessionSecret string `json:"familySessionSecret,omitempty"`
	AccessToken         string `json:"accessToken,omitempty"`
	RefreshToken        string `json:"refreshToken,omitempty"`
}

func PwdLogin(username, password string) (err error) {
	invoker := invoker.New()
	params := url.Values{}
	params.Set("appId", "8025431004")
	params.Set("clientType", "10020")
	params.Set("timeStamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	params.Set("returnURL", "https://m.cloud.189.cn/zhuanti/2020/loginErrorPc/index.html")
	user := &drive.User{Name: username, Password: password}
	resp, err := invoker.PwdLogin("https://cloud.189.cn/unifyLoginForPC.action", params, user)
	if err != nil {
		return err
	}
	var userSession userSession
	params = url.Values{}
	params.Set("version", "6.4.1.0")
	params.Set("clientType", "TELEPC")
	params.Set("channelId", "web_cloud.189.cn")
	params.Set("redirectURL", resp.ToUrl)
	err = invoker.PostForm("https://api.cloud.189.cn/getSessionForPC.action", params, &userSession)
	return
}
