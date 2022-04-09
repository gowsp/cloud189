package invoker

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
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
	refresh func() error
}

func NewInvoker(apiUrl string, refresh func() error, conf *Config) *Invoker {
	jar, _ := cookiejar.New(nil)
	sson := []*http.Cookie{{Name: "SSON", Value: conf.SSON}}
	user := []*http.Cookie{{Name: "COOKIE_LOGIN_USER", Value: conf.Auth}}
	jar.SetCookies(&url.URL{Scheme: "https", Host: "e.189.cn"}, sson)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
	return &Invoker{url: apiUrl, refresh: refresh, http: &http.Client{Jar: jar}, conf: conf}
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

func (i *Invoker) Do(req *http.Request, data interface{}, retry int) error {
	if retry == 0 {
		return os.ErrInvalid
	}
	if i.prepare != nil {
		i.prepare(req)
	}
	resp, err := i.http.Do(req)
	// body, _ := httputil.DumpResponse(resp, true)
	// fmt.Println(string(body))
	if err != nil || resp.StatusCode == http.StatusBadRequest {
		time.Sleep(time.Millisecond * 200)
		err := i.refresh()
		if err != nil {
			return err
		}
		req.Header.Del("Cookie")
		if req.GetBody != nil {
			req.Body, _ = req.GetBody()
		}
		return i.Do(req, data, retry-1)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (i *Invoker) Send(req *http.Request) (*http.Response, error) {
	return i.http.Do(req)
}
func (i *Invoker) Fetch(path string) (*http.Response, error) {
	return i.http.Get(path)
}
func (i *Invoker) Get(path string, params url.Values, data interface{}) error {
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
func (i *Invoker) Post(path string, params url.Values, data interface{}) error {
	url := i.url + path
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return i.Do(req, data, 3)
}
