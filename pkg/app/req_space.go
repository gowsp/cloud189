package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gowsp/cloud189/pkg"
)

func (c *api) Space() (space pkg.Space, err error) {
	err = c.invoker.Get("/getUserInfo.action", nil, &space)
	return
}

type result struct {
	Result    int    `json:"result,omitempty"`
	ResultTip string `json:"resultTip,omitempty"`
}

func (client *api) Sign() error {
	params := url.Values{}
	addParams(&params)
	var r result
	err := client.invoker.Get("/mkt/userSign.action", params, &r)
	if err == nil {
		if r.Result == -1 {
			fmt.Print("已签到 ")
		}
		fmt.Println(r.ResultTip)
	}
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN&activityId=ACT_SIGNIN")
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN_PHOTOS&activityId=ACT_SIGNIN")
	return nil
}

type signResp struct {
	ErrorCode string `json:"errorCode,omitempty"`
	PrizeName string `json:"prizeName,omitempty"`
}

func (a *api) signReq(url string) {
	var e signResp
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	err := a.invoker.Do(req, &e, 3)
	if err == nil {
		switch e.ErrorCode {
		case "User_Not_Chance":
			log.Println("signed")
		case "TimeOut":
			time.Sleep(time.Millisecond * 200)
			a.invoker.Refresh()
			a.signReq(url)
		default:
			log.Printf("obtain: %s", e.PrizeName)
		}
	} else {
		log.Println(err)
	}

}
