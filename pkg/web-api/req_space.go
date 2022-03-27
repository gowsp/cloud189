package web

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gowsp/cloud189/pkg"
)

type space struct {
	Ava uint64 `json:"available,omitempty"`
	Cap uint64 `json:"capacity,omitempty"`
}

func (s *space) Available() uint64 {
	return s.Ava
}

func (s *space) Capacity() uint64 {
	return s.Cap
}

func (c *Api) Space() (pkg.Space, error) {
	var space space
	err := c.invoker.Get("/open/user/getUserInfoForPortal.action", url.Values{}, &space)
	return &space, err
}

type signResp struct {
	ErrorCode string `json:"errorCode,omitempty"`
	PrizeName string `json:"prizeName,omitempty"`
}

func (client *Api) Sign() error {
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN&activityId=ACT_SIGNIN")
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN_PHOTOS&activityId=ACT_SIGNIN")
	return nil
}

func (a *Api) signReq(url string) {
	var e signResp
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	err := a.invoker.Do(req, &e, 3)
	if err == nil {
		if e.ErrorCode == "User_Not_Chance" {
			log.Println("signed")
		} else {
			log.Printf("obtain: %s" + e.PrizeName)
		}
	} else {
		log.Println(err)
	}

}
