package app

import (
	"net/url"
	"strconv"
	"time"

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

func (api *api) PwdLogin(username, password string) (err error) {
	user := &invoker.User{Name: username, Password: password}
	params := url.Values{}
	params.Set("appId", "8025431004")
	params.Set("clientType", "10020")
	params.Set("timeStamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	params.Set("returnURL", "https://m.cloud.189.cn/zhuanti/2020/loginErrorPc/index.html")
	resp, err := api.invoker.PwdLogin("https://cloud.189.cn/unifyLoginForPC.action", params, user)
	if err != nil {
		return err
	}
	var userSession userSession
	params = url.Values{}
	params.Set("redirectURL", resp.ToUrl)
	addParams(&params)
	err = api.invoker.Post("/getSessionForPC.action", params, &userSession)
	api.conf.User = user
	api.conf.Session = &invoker.Session{Key: userSession.SessionKey, Secret: userSession.SessionSecret}
	return api.conf.Save()
}

func addParams(params *url.Values) {
	params.Set("version", "6.4.1.0")
	params.Set("clientType", "TELEPC")
	params.Set("channelId", "web_cloud.189.cn")
}
