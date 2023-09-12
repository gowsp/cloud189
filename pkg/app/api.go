package app

import (
	"fmt"
	"net/http"
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
	return api.PwdLogin(api.conf.User.Name, api.conf.User.Password)
}

func (api *api) sign(req *http.Request) {
	// sha1(SessionKey=相应的值&Operate=相应值&RequestURI=相应值&Date=相应的值”, SessionSecret)
	session := api.conf.Session
	if session == nil {
		return
	}
	date := time.Now().Format(time.RFC1123)
	data := fmt.Sprintf("SessionKey=%s&Operate=%s&RequestURI=%s&Date=%s",
		session.Key, req.Method, req.URL.Path, date)
	req.Header.Set("Date", date)
	req.Header.Set("SessionKey", session.Key)
	req.Header.Set("Signature", util.Sha1(data, session.Secret))
	req.Header.Set("X-Request-ID", util.Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"))
}
