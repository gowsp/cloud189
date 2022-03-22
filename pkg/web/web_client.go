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

	"github.com/gowsp/cloud189/pkg/config"
	"github.com/gowsp/cloud189/pkg/util"
)

var client *Client
var clientSingleton sync.Once

type Client struct {
	api    *http.Client
	config *config.Config
}

// use config
func NewClient(configPath string) *Client {
	config.InitConfigFile(configPath)
	clientSingleton.Do(func() {
		config, err := config.OpenConfig()
		if err != nil {
			config = NewContent().QrLogin()
			config.Save()
		}
		client = &Client{
			config: config,
			api:    newClient(config.SSON, config.Auth),
		}
	})
	return client
}

func NewClientWithUser(name, password string) *Client {
	clientSingleton.Do(func() {
		config := NewContent().PwdLogin(name, password)
		client = &Client{
			config: config,
			api:    newClient(config.SSON, config.Auth),
		}
	})
	return client
}

func newClient(sson, auth string) *http.Client {
	jar, _ := cookiejar.New(nil)
	user := []*http.Cookie{
		{Name: "COOKIE_LOGIN_USER", Value: auth},
	}
	jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
	jar.SetCookies(&url.URL{Scheme: "https", Host: "open.e.189.cn"}, []*http.Cookie{
		{Name: "SSON", Value: sson},
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
	cookie := util.FindCookie(cookies, "COOKIE_LOGIN_USER")
	if cookie != nil {
		config := client.config
		config.Auth = cookie.Value
		config.SessionKey = getUserBriefInfo(*config).SessionKey
		config.Save()
		client.api = newClient(config.SSON, config.Auth)
		return
	}
	user := client.config.User
	if user.Name != "" && user.Password != "" {
		client.config = NewContentWithResp(resp).PwdLogin(user.Name, user.Password)
	} else {
		client.config = NewContentWithResp(resp).QrLogin()
	}
	client.config.Save()
	client.api = newClient(client.config.SSON, client.config.Auth)
}
func (client *Client) rsa() *config.RsaConfig {
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
	client.config.Save()
	return &rsa
}
func (client *Client) initSesstion() {
	user := getUserBriefInfo(*client.config)
	if user.SessionKey == "" {
		client.refresh()
	} else {
		client.config.SessionKey = user.SessionKey
		client.config.Save()
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
		client.config.Save()
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

func getUserBriefInfo(config config.Config) *briefInfo {
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
