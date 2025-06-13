package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/gowsp/cloud189/pkg/util"
)

type api struct {
	invoker *invoker.Invoker
	conf    *invoker.Config
}

func New(path string) *api {
	conf, _ := invoker.OpenConfig(path)
	api := &api{conf: conf}
	api.invoker = invoker.NewInvoker("https://api.cloud.189.cn", api.refresh, conf)
	api.invoker.SetPrepare(api.sign)
	return api
}

func Mem(username, password string) *api {
	conf := &invoker.Config{User: &invoker.User{Name: username, Password: password}}
	api := &api{conf: conf}
	api.invoker = invoker.NewInvoker("https://api.cloud.189.cn", api.refresh, conf)
	api.invoker.SetPrepare(api.sign)
	return api
}

func (api *api) refresh() error {
	user := api.conf.User
	if user.Name == "" || user.Password == "" {
		return errors.New("扫码不支持自动刷新")
	}
	return api.PwdLogin(api.conf.User.Name, api.conf.User.Password)
}

func (api *api) sign(req *http.Request) {
	// sha1(SessionKey=相应的值&Operate=相应值&RequestURI=相应值&Date=相应的值”, SessionSecret)
	session := api.conf.Session
	if session == nil {
		return
	}
	query := req.URL.Query()

	now := time.Now()
	date := now.Format(time.RFC1123)
	data := fmt.Sprintf("SessionKey=%s&Operate=%s&RequestURI=%s&Date=%s",
		session.Key, req.Method, req.URL.Path, date)
	// 追加上传参数
	if req.Host == "upload.cloud.189.cn" {
		data += "&params=" + query.Get("params")
	}
	req.Header.Set("Date", date)
	req.Header.Set("user-agent", "desktop")
	req.Header.Set("SessionKey", session.Key)
	req.Header.Set("Signature", util.Sha1(data, session.Secret))
	req.Header.Set("X-Request-ID", util.Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"))

	// 填充客户端参数
	query.Set("rand", strconv.FormatInt(now.UnixMilli(), 10))
	query.Set("clientType", "TELEPC")
	query.Set("version", "7.1.8.0")
	query.Set("channelId", "web_cloud.189.cn")
	req.URL.RawQuery = query.Encode()
}
