package web

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gowsp/cloud189/pkg/drive"
	"github.com/gowsp/cloud189/pkg/util"
)

type uploader struct {
	api *api
}

type briefInfo struct {
	SessionKey string `json:"sessionKey,omitempty"`
}

func (i *api) session() string {
	if i.sessionKey != "" {
		return i.sessionKey
	}
	var user briefInfo
	i.invoker.Get("/portal/v2/getUserBriefInfo.action", nil, &user)
	i.sessionKey = user.SessionKey
	return i.sessionKey
}
func (i *api) rsa() *drive.RsaConfig {
	rsa := i.conf.RSA
	if rsa.Expire > time.Now().UnixMilli() {
		return &rsa
	}
	i.invoker.Get("/security/generateRsaKey.action", nil, &rsa)
	i.conf.RSA = rsa
	i.conf.Save()
	return &rsa
}

type uploadResp interface {
	GetCode() string
}

func (uploader *api) do(u string, f url.Values, result uploadResp) error {
	c := strconv.FormatInt(time.Now().UnixMilli(), 10)
	r := util.Random("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	l := util.Random("xxxxxxxxxxxx4xxxyxxxxxxxxxxxxxxx")
	l = l[0 : 16+int(16*rand.Float32())]

	e := util.EncodeParam(f)
	data := util.AesEncrypt([]byte(e), []byte(l[0:16]))
	h := hex.EncodeToString(data)

	req, err := http.NewRequest(http.MethodGet, "https://upload.cloud.189.cn"+u+"?params="+h, nil)
	if err != nil {
		return err
	}
	a := make(url.Values)
	a.Set("SessionKey", uploader.session())
	a.Set("Operate", http.MethodGet)
	a.Set("RequestURI", u)
	a.Set("Date", c)
	a.Set("params", h)

	req.Header.Set("accept", "application/json;charset=UTF-8")
	req.Header.Set("SessionKey", uploader.session())

	g := util.SHA1(util.EncodeParam(a), l)
	req.Header.Set("Signature", hex.EncodeToString(g))
	req.Header.Set("X-Request-Date", c)
	req.Header.Set("X-Request-ID", r)

	b := uploader.rsa().Encrypt(l)
	req.Header.Set("EncryptionText", base64.StdEncoding.EncodeToString(b))
	req.Header.Set("PkId", uploader.rsa().PkId)

	if err = uploader.invoker.Do(req, result, 3); err != nil {
		return err
	}
	switch result.GetCode() {
	case "SUCCESS":
		return nil
	case "InvalidSessionKey":
		uploader.refresh()
		uploader.sessionKey = ""
		return uploader.do(u, f, result)
	case "InvalidSignature":
		uploader.refresh()
		uploader.sessionKey = ""
		return uploader.do(u, f, result)
	default:
		return errors.New(result.GetCode())
	}
}
