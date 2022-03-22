package web

import (
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/gowsp/cloud189/pkg/util"
)

func NewContent() *Content {
	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/loginUrl.action", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return NewContentWithResp(resp)
}

func NewContentWithResp(resp *http.Response) *Content {
	data, _ := io.ReadAll(resp.Body)
	location := resp.Request.Response.Header.Get("location")
	cookie := util.FindCookie(resp.Cookies(), "LT")
	content := Content{Cookie: cookie, data: data, Referer: location}
	content.parse()
	return &content
}

type Content struct {
	Cookie      *http.Cookie
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

func (c *Content) parse() {
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

func (c *Content) read(str string) string {
	reg := regexp.MustCompile(str)
	paramId := reg.FindSubmatch(c.data)
	return string(paramId[1])
}
