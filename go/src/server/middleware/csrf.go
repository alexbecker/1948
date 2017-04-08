package middleware

import (
	"net/url"
	"os"
	"strings"
)

var (
	hosts []string
)

func init() {
	hosts = strings.Split(os.Getenv("HOSTS"), ",")
}

func IsSafeMethod(method string) bool {
	return (method == "GET" || method == "HEAD" || method == "OPTIONS")
}

func CheckReferer(referer string) bool {
	if referer == "" {
		return true
	}

	parsed, err := url.ParseRequestURI(referer)
	if err != nil {
		return false
	}
	for i := 0; i < len(hosts); i++ {
		if hosts[i] == parsed.Host {
			return true
		}
	}
	return false
}
