package util

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func Debug(req *http.Request) *http.Response {
	data, _ := httputil.DumpRequest(req, true)
	log.Println(string(data))
	resp, _ := http.DefaultClient.Do(req)
	data, _ = httputil.DumpResponse(resp, true)
	log.Println(string(data))
	return resp
}
