package invoker

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gowsp/cloud189/pkg/util"
)

type content struct {
	http    *http.Client
	user    *User
	Referer string
	AppKey  string
	ReqId   string
	Lt      string
}

func newLoginContent(http *http.Client, user *User, Referer string) *content {
	v, _ := url.Parse(Referer)
	lt := v.Query().Get("lt")
	reqId := v.Query().Get("reqId")
	appKey := v.Query().Get("appId")
	return &content{http: http, user: user, AppKey: appKey, Referer: Referer, Lt: lt, ReqId: reqId}
}

type appConf struct {
	Data struct {
		AccountType          string `json:"accountType"`
		AgreementCheck       string `json:"agreementCheck"`
		AppKey               string `json:"appKey"`
		ClientType           int    `json:"clientType"`
		DefaultSaveName      string `json:"defaultSaveName"`
		DefaultSaveNameCheck string `json:"defaultSaveNameCheck"`
		IsOauth2             bool   `json:"isOauth2"`
		LoginSort            string `json:"loginSort"`
		MailSuffix           string `json:"mailSuffix"`
		PageKey              string `json:"pageKey"`
		ParamID              string `json:"paramId"`
		RegReturnURL         string `json:"regReturnUrl"`
		ReqID                string `json:"reqId"`
		ReturnURL            string `json:"returnUrl"`
		ShowFeedback         string `json:"showFeedback"`
		ShowPwSaveName       string `json:"showPwSaveName"`
		ShowQrSaveName       string `json:"showQrSaveName"`
		ShowSmsSaveName      string `json:"showSmsSaveName"`
		Sso                  string `json:"sso"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Result string `json:"result"`
}

func (ctx *content) getAppConf() *appConf {
	params := make(url.Values)
	params.Set("version", "2.0")
	params.Set("appKey", ctx.AppKey)
	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/appConf.do",
		strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://open.e.189.cn")
	req.Header.Set("Referer", ctx.Referer)
	req.Header.Set("Reqid", ctx.ReqId)
	req.Header.Set("lt", ctx.Lt)
	resp, err := ctx.http.Do(req)
	if err != nil {
		return nil
	}
	var appConf appConf
	json.NewDecoder(resp.Body).Decode(&appConf)
	return &appConf
}

type encryptConf struct {
	Result int `json:"result"`
	Data   struct {
		UpSmsOn   string `json:"upSmsOn"`
		Pre       string `json:"pre"`
		PreDomain string `json:"preDomain"`
		PubKey    string `json:"pubKey"`
	} `json:"data"`
}

func (ctx *content) getEncryptConf() *encryptConf {
	params := make(url.Values)
	params.Set("appId", "cloud")
	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/config/encryptConf.do",
		strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ctx.Referer)
	resp, err := ctx.http.Do(req)
	if err != nil {
		return nil
	}
	var encryptConf encryptConf
	json.NewDecoder(resp.Body).Decode(&encryptConf)
	return &encryptConf
}
func (ctx *content) pwdRequest() *http.Request {
	appConf := ctx.getAppConf().Data
	encryptConf := ctx.getEncryptConf().Data
	user := ctx.user
	key := util.Key(encryptConf.PubKey)
	data, _ := util.RsaEncrypt(key, []byte(user.Name))
	name := hex.EncodeToString(data)
	data, _ = util.RsaEncrypt(key, []byte(user.Password))
	password := hex.EncodeToString(data)

	params := make(url.Values)
	params.Set("version", "v2.0")
	params.Set("appKey", appConf.AppKey)
	params.Set("accountType", appConf.AccountType)
	params.Set("userName", encryptConf.Pre+name)
	params.Set("epd", encryptConf.Pre+password)
	params.Set("captchaType", "")
	params.Set("validateCode", "")
	params.Set("smsValidateCode", "")
	params.Set("captchaToken", "")
	params.Set("returnUrl", ctx.Referer)
	params.Set("mailSuffix", appConf.MailSuffix)
	params.Set("dynamicCheck", "FALSE")
	params.Set("clientType", strconv.Itoa(appConf.ClientType))
	params.Set("cb_SaveName", "0")
	params.Set("isOauth2", strconv.FormatBool(appConf.IsOauth2))
	params.Set("state", "")
	params.Set("paramId", appConf.ParamID)

	req, _ := http.NewRequest(http.MethodPost, "https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do",
		strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ctx.Referer)
	req.Header.Set("Reqid", ctx.ReqId)
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
	location := resp.Request.Response.Header.Get("location")
	content := newLoginContent(i.http, user, location)
	return content, nil
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
