package invoker

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gowsp/cloud189/pkg/util"
)

type content struct {
	user        *User
	Referer     string
	Captcha     string
	RsaKey      string
	data        []byte
	AppKey      string
	ReqId       string
	IsOauth2    string
	ParamId     string
	ReturnUrl   string
	ClientType  string
	AccountType string
	MailSuffix  string
	Lt          string
}

func (c *content) parse() {
	c.AppKey = c.read("appKey = '(\\w+)'")
	c.ReqId = c.read("reqId = \"(\\w+)\"")
	c.RsaKey = c.read("\"j_rsaKey\" value=\"(.+)\"")
	c.ParamId = c.read("paramId = \"(\\w+)\"")
	c.IsOauth2 = c.read("isOauth2 = \"(\\w+)\"")
	c.ClientType = c.read("clientType = '(\\w+)'")
	c.AccountType = c.read("accountType = '(\\w+)'")
	c.MailSuffix = c.read("mailSuffix = '(.*?)'")
	c.Captcha = c.read("'captchaToken' value='(.+)'")
	c.ReturnUrl = c.read("returnUrl = '(.+)'")
	c.Lt = c.read("lt = \"(\\w+)\"")
}
func (c *content) read(str string) string {
	reg := regexp.MustCompile(str)
	paramId := reg.FindSubmatch(c.data)
	return string(paramId[1])
}

func (ctx *content) pwdRequest() *http.Request {
	user := ctx.user
	key := util.Key(ctx.RsaKey)
	data, _ := util.RsaEncrypt(key, []byte(user.Name))
	name := hex.EncodeToString(data)
	data, _ = util.RsaEncrypt(key, []byte(user.Password))
	password := hex.EncodeToString(data)

	params := make(url.Values)
	params.Set("appKey", ctx.AppKey)
	params.Set("accountType", ctx.AccountType)
	params.Set("userName", "{RSA}"+name)
	params.Set("password", "{RSA}"+password)
	params.Set("validateCode", "")
	params.Set("smsValidateCode", "")
	params.Set("captchaToken", ctx.Captcha)
	params.Set("returnUrl", ctx.Referer)
	params.Set("mailSuffix", ctx.MailSuffix)
	params.Set("dynamicCheck", "FALSE")
	params.Set("clientType", ctx.ClientType)
	params.Set("cb_SaveName", "0")
	params.Set("isOauth2", ctx.IsOauth2)
	params.Set("state", "")
	params.Set("paramId", ctx.ParamId)

	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do",
		strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ctx.Referer)
	req.Header.Set("reqid", ctx.ReqId)
	req.Header.Set("lt", ctx.Lt)
	return req
}

type LoginResult struct {
	Result int    `json:"result,omitempty"`
	Msg    string `json:"msg,omitempty"`
	ToUrl  string `json:"toUrl,omitempty"`
	SSON   string
}

func (i *Invoker) prepareLogin(link string, params url.Values, user *User) (result *content, err error) {
	req, err := util.GetReq(link, params)
	if err != nil {
		return nil, err
	}
	resp, err := i.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	location := resp.Request.Response.Header.Get("location")
	content := content{user: user, data: data, Referer: location}
	content.parse()
	return &content, nil
}

func (i *Invoker) PwdLogin(link string, params url.Values, user *User) (result *LoginResult, err error) {
	content, err := i.prepareLogin(link, params, user)
	if err != nil {
		return nil, err
	}
	resp, err := i.http.Do(content.pwdRequest())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result = &LoginResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return
	}
	if result.Result != 0 {
		return nil, errors.New(result.Msg)
	}
	result.SSON = util.FindCookieValue(resp.Cookies(), "SSON")
	i.conf.SSON = result.SSON
	return
}
