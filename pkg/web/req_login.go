package web

import (
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/util"
)

func login(resp *http.Response, user drive.User) *content {
	data, _ := io.ReadAll(resp.Body)
	location := resp.Request.Response.Header.Get("location")
	content := content{user: user, data: data, Referer: location}
	content.parse()
	return &content
}

type content struct {
	user        drive.User
	Referer     string
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
	c.MailSuffix = c.read("mailSuffix = '(.+)'")
	c.ReturnUrl = c.read("returnUrl = '(.+)'")
	c.Lt = c.read("lt = \"(\\w+)\"")
}
func (c *content) read(str string) string {
	reg := regexp.MustCompile(str)
	paramId := reg.FindSubmatch(c.data)
	return string(paramId[1])
}
func (ctx *content) toRequest() *http.Request {
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
	params.Set("returnUrl", ctx.Referer)
	params.Set("mailSuffix", ctx.MailSuffix)
	params.Set("dynamicCheck", "FALSE")
	params.Set("clientType", ctx.ClientType)
	params.Set("isOauth2", ctx.IsOauth2)
	params.Set("state", "")
	params.Set("paramId", ctx.ParamId)

	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do", strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ctx.Referer)
	req.Header.Set("reqid", ctx.ReqId)
	req.Header.Set("lt", ctx.Lt)
	return req
}

type pwdLoginResult struct {
	Result int    `json:"result,omitempty"`
	ToUrl  string `json:"toUrl,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

func (c *api) Login(name, password string) error {
	user := drive.User{Name: name, Password: password}
	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/loginUrl.action", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	ctx := login(resp, user)
	return c.invoker.Login(ctx)
}
