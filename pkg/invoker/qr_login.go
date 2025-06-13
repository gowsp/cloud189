package invoker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/util"
)

type QrCodeReq struct {
	content    *content
	Uuid       string `json:"uuid,omitempty"`
	Encryuuid  string `json:"encryuuid,omitempty"`
	Encodeuuid string `json:"encodeuuid,omitempty"`
}

func (c *content) qrLogin() {
	config := c.getAppConf()

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
		status := ctx.query(config)
		switch status.Status {
		case -106:
			log.Println("not scanned")
		case -11002:
			log.Println("unconfirmed")
		case 0:
			t.Stop()
			log.Println("logged")
		default:
			t.Stop()
			log.Fatalln("unknown status")
		}
		<-t.C
	}
}

type qrCodeState struct {
	RedirectUrl string `json:"redirectUrl,omitempty"`
	Status      int32  `json:"status,omitempty"`
	SSON        string
}

func (c *QrCodeReq) query(conf *appConf) qrCodeState {
	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/qrcodeLoginState.do", nil)
	req.Header.Set("referer", c.content.Referer)
	params := req.URL.Query()
	params.Set("appId", conf.Data.AppKey)
	params.Set("encryuuid", c.Encryuuid)
	params.Set("date", time.Now().Format("2006-01-0215:04:059"))
	params.Set("uuid", c.Uuid)
	params.Set("returnUrl", conf.Data.ReturnURL)
	params.Set("clientType", strconv.Itoa(conf.Data.ClientType))
	params.Set("timeStamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	params.Set("cb_SaveName", "0")
	params.Set("isOauth2", strconv.FormatBool(conf.Data.IsOauth2))
	params.Set("state", "")
	params.Set("paramId", conf.Data.ParamID)
	req.URL.RawQuery = params.Encode()

	resp, _ := http.DefaultClient.Do(req)
	var status qrCodeState

	json.NewDecoder(resp.Body).Decode(&status)

	if status.Status != 0 {
		return status
	}
	status.SSON = util.FindCookie(resp.Cookies(), "SSON").Value
	return status
}

func (i *Invoker) QrLogin(link string, params url.Values) (result *LoginResult, err error) {
	content, err := i.prepareLogin(link, params, nil)
	if err != nil {
		return nil, err
	}
	config := content.getAppConf()
	req, _ := http.NewRequest(http.MethodGet, "https://open.e.189.cn/api/logbox/oauth2/getUUID.do", nil)
	param := req.URL.Query()
	param.Set("appId", content.AppKey)
	req.URL.RawQuery = param.Encode()
	resp, _ := http.DefaultClient.Do(req)
	var ctx QrCodeReq
	ctx.content = content
	json.NewDecoder(resp.Body).Decode(&ctx)
	params = make(url.Values)
	url, _ := url.PathUnescape(ctx.Encodeuuid)
	params.Set("REQID", content.ReqId)
	params.Set("uuid", url)
	log.Printf("please open url in your browser to login:\nhttps://open.e.189.cn/api/logbox/oauth2/image.do?%s\n\n", params.Encode())
	t := time.NewTicker(3 * time.Second)
	var status qrCodeState
	for {
		status = ctx.query(config)
		switch status.Status {
		case -106:
			log.Println("not scanned")
		case -11002:
			log.Println("unconfirmed")
		case 0:
			t.Stop()
			log.Println("logged")
			return &LoginResult{ToUrl: status.RedirectUrl, SSON: status.SSON}, nil
		default:
			t.Stop()
			return nil, errors.New("unknown status")
		}
		<-t.C
	}
}
