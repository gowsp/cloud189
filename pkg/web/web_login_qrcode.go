package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/config"
	"github.com/gowsp/cloud189/pkg/util"
)

type QrCodeReq struct {
	content    *Content
	Uuid       string `json:"uuid,omitempty"`
	Encryuuid  string `json:"encryuuid,omitempty"`
	Encodeuuid string `json:"encodeuuid,omitempty"`
}

func (c *Content) QrLogin() *config.Config {
	req, _ := http.NewRequest(http.MethodGet, "https://open.e.189.cn/api/logbox/oauth2/getUUID.do", nil)
	param := req.URL.Query()
	param.Set("appId", c.AppKey)
	req.URL.RawQuery = param.Encode()
	resp, _ := http.DefaultClient.Do(req)
	var ctx QrCodeReq
	ctx.content = c
	json.NewDecoder(resp.Body).Decode(&ctx)
	params := make(url.Values)
	url, _ := url.PathUnescape(ctx.Encodeuuid)
	params.Set("REQID", c.ReqId)
	params.Set("uuid", url)
	log.Printf("please open url in your browser to login:\nhttps://open.e.189.cn/api/logbox/oauth2/image.do?%s\n\n", params.Encode())
	t := time.NewTicker(3 * time.Second)
	for {
		status := ctx.query()
		switch status.Status {
		case -106:
			log.Println("not scanned")
		case -11002:
			log.Println("unconfirmed")
		case 0:
			t.Stop()
			log.Println("logged")
			config := config.Config{}
			config.SsonLogin(status.RedirectUrl, status.SSON)
			return &config
		default:
			log.Fatalln("unknown status")
		}
		<-t.C
	}
}

type qrCodeState struct {
	RedirectUrl string `json:"redirectUrl,omitempty"`
	Status      int32  `json:"status,omitempty"`
	SSON        *http.Cookie
}

func (c *QrCodeReq) query() qrCodeState {
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

	resp, _ := http.DefaultClient.Do(req)
	var status qrCodeState

	json.NewDecoder(resp.Body).Decode(&status)

	if status.Status != 0 {
		return status
	}
	status.SSON = util.FindCookie(resp.Cookies(), "SSON")
	return status
}
