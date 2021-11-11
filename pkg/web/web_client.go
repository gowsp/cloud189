package web

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

var client *Client
var clientSingleton sync.Once

type Client struct {
	api    *http.Client
	config *Config
}

func GetClient() *Client {
	clientSingleton.Do(func() {
		config := GetConfig()
		client = &Client{
			config: config,
			api:    newClient(config.SSON, config.Auth),
		}
	})
	return client
}

func EmptyClient() *Client {
	return &Client{}
}

func newClient(sson, auth string) *http.Client {
	jar, _ := cookiejar.New(nil)
	user := []*http.Cookie{
		{Name: "COOKIE_LOGIN_USER", Value: config.Auth},
	}
	jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "open.e.189.cn"}, []*http.Cookie{
		{Name: "SSON", Value: config.SSON},
	})
	return &http.Client{Jar: jar}
}

func (client *Client) refresh() {
	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/loginUrl.action", nil)

	resp, err := client.api.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	cookies := client.api.Jar.Cookies(resp.Request.URL)
	cookie := findCookie(cookies, "COOKIE_LOGIN_USER")
	if cookie != nil {
		config := client.config
		config.Auth = cookie.Value
		config.SessionKey = getUserBriefInfo(*config).SessionKey
		config.save()
		client.api = newClient(config.SSON, config.Auth)
		return
	}
	user := config.User
	if user.Name != "" && user.Password != "" {
		NewContentWithResp(resp).PwdLogin(user.Name, user.Password)
	} else {
		NewContentWithResp(resp).QrLogin()
	}
	client.api = newClient(config.SSON, config.Auth)
}
func (client *Client) rsa() *rsa {
	config := client.config
	rsa := client.config.RSA
	now := time.Now().UnixMilli()
	if rsa.Expire > now {
		return &rsa
	}
	for rsa.Expire < now {
		req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/security/generateRsaKey.action", nil)
		req.Header.Add("accept", "application/json;charset=UTF-8")
		resp, _ := client.api.Do(req)
		json.NewDecoder(resp.Body).Decode(&rsa)
		resp.Body.Close()
		if resp.StatusCode != 200 || rsa.Expire == 0 {
			client.refresh()
		}
	}
	config.RSA = rsa
	config.save()
	return &rsa
}
func (client *Client) initSesstion() {
	user := getUserBriefInfo(*client.config)
	if user.SessionKey == "" {
		client.refresh()
	} else {
		config.SessionKey = user.SessionKey
		config.save()
	}
}
func (client *Client) sesstionKey() string {
	config := client.config
	key := config.SessionKey
	if key != "" {
		return key
	}
	user := getUserBriefInfo(*client.config)
	if user.SessionKey == "" {
		client.refresh()
	} else {
		config.SessionKey = user.SessionKey
		config.save()
	}
	return config.SessionKey
}

type errorResp struct {
	ErrorCode string `json:"errorCode,omitempty"`
}

func (resp *errorResp) IsInvalidSession() bool {
	return resp.ErrorCode == "InvalidSessionKey"
}

type briefInfo struct {
	SessionKey  string `json:"sessionKey,omitempty"`
	UserAccount string `json:"userAccount,omitempty"`
}

func getUserBriefInfo(config Config) *briefInfo {
	u := fmt.Sprintf("https://cloud.189.cn/v2/getUserBriefInfo.action?noCache=%v", rand.Float64())
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.AddCookie(&http.Cookie{Name: "COOKIE_LOGIN_USER", Value: config.Auth})
	req.Header.Add("accept", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var user briefInfo
	if strings.Index(resp.Header.Get("Content-Type"), "html") > 0 {
		return &user
	}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user
}

func (client *Client) isInvalidSession(data []byte) (invalid bool) {
	var errorResp errorResp
	json.Unmarshal(data, &errorResp)
	invalid = errorResp.IsInvalidSession()
	if invalid {
		client.refresh()
	}
	return
}
