package middleware

import (
	"net/http"
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

// Checks that the origin matches one of the hosts specified in HOSTS.
// Falls back to Referer if the Origin header is not present.
// Accepts if neither header is present.
func CheckOrigin(req *http.Request) bool {
	origin := ""
	if len(req.Header["Origin"]) > 0 {
		origin = req.Header["Origin"][0]
	}

	if origin == "" {
		origin = req.Referer()
		if origin == "" {
			return true
		}
	}

	parsed, err := url.ParseRequestURI(origin)
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
