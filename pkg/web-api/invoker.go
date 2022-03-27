package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/util"
)

const domain = "https://cloud.189.cn/api"

type invoker struct {
	http *http.Client
	conf *drive.Config
}

const COOKIE_USER = "COOKIE_LOGIN_USER"

func newInvoker(conf *drive.Config) *invoker {
	jar, _ := cookiejar.New(nil)
	sson := []*http.Cookie{{Name: "SSON", Value: conf.SSON}}
	user := []*http.Cookie{{Name: COOKIE_USER, Value: conf.Auth}}
	jar.SetCookies(&url.URL{Scheme: "https", Host: "e.189.cn"}, sson)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
	return &invoker{http: &http.Client{Jar: jar}, conf: conf}
}

func (i *invoker) refresh() error {
	req, _ := http.NewRequest(http.MethodGet, domain+"/portal/loginUrl.action", nil)

	resp, err := i.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	cookies := i.http.Jar.Cookies(resp.Request.URL)
	user := util.FindCookie(cookies, COOKIE_USER)
	if user != nil {
		i.conf.Auth = user.Value
		i.conf.Save()
		return nil
	}
	return i.Login(LoginContent(resp, i.conf.User))
}

func (i *invoker) Do(req *http.Request, data interface{}, retry int) error {
	if retry == 0 {
		return os.ErrInvalid
	}
	resp, err := i.http.Do(req)
	body, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(body))
	if err != nil || resp.StatusCode == http.StatusBadRequest {
		time.Sleep(time.Millisecond * 200)
		i.refresh()
		req.Header.Del("Cookie")
		if req.GetBody != nil {
			req.Body, _ = req.GetBody()
		}
		return i.Do(req, data, retry-1)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}
func (i *invoker) Get(path string, params url.Values, data interface{}) error {
	url := domain + path
	if len(params) > 0 {
		url += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	return i.Do(req, data, 3)
}
func (i *invoker) Post(path string, params url.Values, data interface{}) error {
	url := domain + path
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return i.Do(req, data, 3)
}
func (i *invoker) Login(ctx *Content) error {
	req := ctx.toRequest()
	resp, _ := i.http.Do(req)
	var result pwdLoginResult
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	if result.Result != 0 {
		return fmt.Errorf("login failed")
	}
	sson := util.FindCookie(resp.Cookies(), "SSON")
	resp, _ = i.http.Get(result.ToUrl)
	user := util.FindCookie(i.http.Jar.Cookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}), "COOKIE_LOGIN_USER")
	i.conf.User = ctx.user
	i.conf.SSON = sson.Value
	i.conf.Auth = user.Value
	i.conf.Save()
	return nil
}
