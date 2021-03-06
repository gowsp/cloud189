package util

import "net/http"

// Find special name cookie
func FindCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookies := range cookies {
		if cookies.Name == name {
			return cookies
		}
	}
	return nil
}
func FindCookieValue(cookies []*http.Cookie, name string) string {
	for _, cookies := range cookies {
		if cookies.Name == name {
			return cookies.Value
		}
	}
	return ""
}
