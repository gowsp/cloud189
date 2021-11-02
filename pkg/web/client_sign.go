package web

import (
	"encoding/json"
	"log"
)

type signResp struct {
	ErrorCode string `json:"errorCode,omitempty"`
	PrizeName string `json:"prizeName,omitempty"`
}

func (client *Client) Sign() {
	client.initSesstion()
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN&activityId=ACT_SIGNIN")
	client.signReq("https://m.cloud.189.cn/v2/drawPrizeMarketDetails.action?taskId=TASK_SIGNIN_PHOTOS&activityId=ACT_SIGNIN")
}

func (client *Client) signReq(url string) {
	var e signResp
	resp, err := client.api.Get(url)
	if err == nil {
		json.NewDecoder(resp.Body).Decode(&e)
		if e.ErrorCode == "User_Not_Chance" {
			log.Println("signed")
		} else {
			log.Printf("obtain: %s" + e.PrizeName)
		}
	} else {
		log.Println(err)
	}

}
