package invoker

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
)

type Invoker struct {
	url     string
	http    *http.Client
	conf    *Config
	prepare func(*http.Request)
	Refresh func() error
}

func NewInvoker(apiUrl string, refresh func() error, conf *Config) *Invoker {
	jar, _ := cookiejar.New(nil)
	sson := []*http.Cookie{{Name: "SSON", Value: conf.SSON}}
	user := []*http.Cookie{{Name: "COOKIE_LOGIN_USER", Value: conf.Auth}}
	jar.SetCookies(&url.URL{Scheme: "https", Host: "e.189.cn"}, sson)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
	return &Invoker{url: apiUrl, Refresh: refresh, http: &http.Client{Jar: jar}, conf: conf}
}

func (i *Invoker) SetPrepare(prepare func(req *http.Request)) {
	i.prepare = prepare
}
func (i *Invoker) Cookies(url *url.URL) []*http.Cookie {
	return i.http.Jar.Cookies(url)
}
func (i *Invoker) Cookie(raw, name string) string {
	url, _ := url.Parse(raw)
	cookies := i.http.Jar.Cookies(url)
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func (i *Invoker) DoWithResp(req *http.Request) (*http.Response, error) {
	if i.prepare != nil {
		i.prepare(req)
	}
	resp, err := i.http.Do(req)
	val := os.Getenv("189_MODE")
	if val == "1" {
		rdata, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(rdata))
		data, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(data))
	}
	return resp, err
}
func (i *Invoker) Do(req *http.Request, data any, retry int) error {
	if retry == 0 {
		return os.ErrInvalid
	}
	resp, err := i.DoWithResp(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			return err
		}
		if rsp, ok := data.(OkRsp); ok && !rsp.IsSuccess() {
			return rsp
		}
		return nil

	case http.StatusBadRequest, http.StatusForbidden:
		rsp := new(strCodeRsp)
		err := json.NewDecoder(resp.Body).Decode(rsp)
		if err == nil && rsp.isBusinessErr() {
			return rsp
		}
		time.Sleep(time.Millisecond * 200)
		err = i.Refresh()
		if err != nil {
			return err
		}
		req.Header.Del("Cookie")
		if req.GetBody != nil {
			req.Body, _ = req.GetBody()
		}
		return i.Do(req, data, retry-1)
	default:
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
}

func (i *Invoker) Send(req *http.Request) (*http.Response, error) {
	return i.http.Do(req)
}
func (i *Invoker) Fetch(path string) (*http.Response, error) {
	return i.http.Get(path)
}
func (i *Invoker) Get(path string, params url.Values, data any) error {
	url := i.url + path
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
func (i *Invoker) Post(path string, params url.Values, data any) error {
	url := i.url + path
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return i.Do(req, data, 3)
}
