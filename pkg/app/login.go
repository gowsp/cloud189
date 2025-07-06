package app

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/invoker"
)

func (api *api) beforLogin() url.Values {
	params := url.Values{}
	params.Set("appId", "9317140619")
	params.Set("clientType", "10020")
	params.Set("timeStamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	params.Set("returnURL", "https://m.cloud.189.cn/zhuanti/2020/loginErrorPc/index.html")
	return params
}

func (api *api) PwdLogin(username, password string) (err error) {
	params := api.beforLogin()
	user := &invoker.User{Name: username, Password: password}
	resp, err := api.invoker.PwdLogin("https://cloud.189.cn/unifyLoginForPC.action", params, user)
	if err != nil {
		return err
	}
	api.conf.User = user
	return api.afterLogin(resp)
}

func (api *api) QrLogin() (err error) {
	params := api.beforLogin()
	resp, err := api.invoker.QrLogin("https://cloud.189.cn/unifyLoginForPC.action", params)
	if err != nil {
		return err
	}
	return api.afterLogin(resp)
}

func (api *api) afterLogin(resp *invoker.LoginResult) error {
	var userSession invoker.Session
	params := url.Values{}
	params.Set("redirectURL", resp.ToUrl)
	if err := api.invoker.Post("/getSessionForPC.action", params, &userSession); err != nil {
		return err
	}
	api.conf.Session = &userSession
	return api.conf.Save()
}
