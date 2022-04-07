package util

import (
	"net/http"
	"net/url"
	"strings"
)

func GetReq(u string, params url.Values) (*http.Request, error) {
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	return http.NewRequest(http.MethodGet, u, nil)
}

func EncodeParam(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	for _, k := range keys {
		vs := v[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}
