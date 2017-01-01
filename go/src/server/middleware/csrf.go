package middleware

import (
	"net/url"
	"os"
)

var (
	host string
)

func init() {
	host = os.Getenv("HOST")
}

func IsSafeMethod(method string) bool {
	return (method == "GET" || method == "HEAD" || method == "OPTIONS")
}

func CheckReferer(referer string) bool {
	if referer == "" {
		return true
	}

	parsed, err := url.ParseRequestURI(referer)
	return (err == nil && parsed.Host == host)
}
