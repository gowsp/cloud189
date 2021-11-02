package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

func (c *Client) Login() {
	NewContent().QrCode().Login()
}

type QrCodeReq struct {
	content    *Content
	client     *http.Client
	Uuid       string `json:"uuid,omitempty"`
	Encryuuid  string `json:"encryuuid,omitempty"`
	Encodeuuid string `json:"encodeuuid,omitempty"`
}

func (c *Content) QrCode() *QrCodeReq {
	req, _ := http.NewRequest(http.MethodGet, "https://open.e.189.cn/api/logbox/oauth2/getUUID.do", nil)
	param := req.URL.Query()
	param.Set("appId", c.AppKey)
	req.URL.RawQuery = param.Encode()
	resp, _ := http.DefaultClient.Do(req)
	var ctx QrCodeReq
	json.NewDecoder(resp.Body).Decode(&ctx)
	ctx.content = c
	jar, _ := cookiejar.New(nil)
	ctx.client = &http.Client{Jar: jar}
	return &ctx
}

func (c *QrCodeReq) Login() {
	params := make(url.Values)
	url, _ := url.PathUnescape(c.Encodeuuid)
	params.Set("REQID", c.content.ReqId)
	params.Set("uuid", url)
	log.Printf("please open url in your browser to login:\nhttps://open.e.189.cn/api/logbox/oauth2/image.do?%s\n\n", params.Encode())
	t := time.NewTicker(3 * time.Second)
	for {
		status := c.query()
		switch status.Status {
		case -106:
			log.Println("not scanned")
		case -11002:
			log.Println("unconfirmed")
		case 0:
			t.Stop()
			log.Println("logged")
			resp, _ := c.client.Get(status.RedirectUrl)
			cookies := c.client.Jar.Cookies(resp.Request.URL)
			config = Config{SSON: status.SSON}
			for _, cookie := range cookies {
				if cookie.Name == "COOKIE_LOGIN_USER" {
					config.Auth = cookie.Value
					config.Save()
					return
				}
			}
			return
		default:
			log.Fatalln("unknown status")
		}
		<-t.C
	}
}

type QrCodeState struct {
	RedirectUrl string `json:"redirectUrl,omitempty"`
	Status      int32  `json:"status,omitempty"`
	SSON        string
}

func (c *QrCodeReq) query() QrCodeState {
	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/qrcodeLoginState.do", nil)
	req.Header.Set("referer", c.content.Referer)
	params := req.URL.Query()
	params.Set("appId", c.content.AppKey)
	params.Set("encryuuid", c.Encryuuid)
	params.Set("date", time.Now().Format("2006-01-0215:04:059"))
	params.Set("uuid", c.Uuid)
	params.Set("returnUrl", c.content.ReturnUrl)
	params.Set("clientType", c.content.ClientType)
	params.Set("timeStamp", strconv.FormatInt(time.Now().Unix()*1000, 10))
	params.Set("cb_SaveName", "0")
	params.Set("isOauth2", c.content.IsOauth2)
	params.Set("state", "")
	params.Set("paramId", c.content.ParamId)
	req.URL.RawQuery = params.Encode()

	resp, _ := c.client.Do(req)
	var status QrCodeState

	json.NewDecoder(resp.Body).Decode(&status)

	if status.Status != 0 {
		return status
	}
	cookies := c.client.Jar.Cookies(resp.Request.URL)
	for _, cookie := range cookies {
		if cookie.Name == "SSON" {
			status.SSON = cookie.Value
			break
		}
	}
	return status
}
