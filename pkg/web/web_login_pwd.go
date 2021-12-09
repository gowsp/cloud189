package web

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gowsp/cloud189-cli/pkg/config"
	"github.com/gowsp/cloud189-cli/pkg/util"
)

type pwdLoginResult struct {
	Result int    `json:"result,omitempty"`
	ToUrl  string `json:"toUrl,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

func (ctx *Content) PwdLogin(name, password string) *config.Config {
	user := config.User{Name: name, Password: password}
	key := util.Key(ctx.RsaKey)
	data, _ := util.RsaEncrypt(key, []byte(name))
	name = hex.EncodeToString(data)
	data, _ = util.RsaEncrypt(key, []byte(password))
	password = hex.EncodeToString(data)

	params := make(url.Values)
	params.Set("appKey", ctx.AppKey)
	params.Set("accountType", ctx.AccountType)
	params.Set("userName", "{RSA}"+name)
	params.Set("password", "{RSA}"+password)
	params.Set("validateCode", "")
	params.Set("returnUrl", ctx.Referer)
	params.Set("mailSuffix", ctx.MailSuffix)
	params.Set("dynamicCheck", "FALSE")
	params.Set("clientType", ctx.ClientType)
	params.Set("isOauth2", ctx.IsOauth2)
	params.Set("state", "")
	params.Set("paramId", ctx.ParamId)

	req, _ := http.NewRequest("POST", "https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do", strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ctx.Referer)
	req.Header.Set("reqid", ctx.ReqId)
	req.Header.Set("lt", ctx.Lt)
	req.AddCookie(ctx.Cookie)

	resp, _ := http.DefaultClient.Do(req)
	var result pwdLoginResult
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	if result.Result == 0 {
		sson := util.FindCookie(resp.Cookies(), "SSON")
		config := config.Config{User: user}
		config.SsonLogin(result.ToUrl, sson)
		return &config
	}
	fmt.Println(result.Msg)
	return nil
}
