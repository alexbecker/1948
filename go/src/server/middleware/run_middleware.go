package middleware

import (
	"net/http"
)

func MiddlewareServeMux(mux *http.ServeMux) *http.ServeMux {
	newMux := http.NewServeMux()
	newMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		wLogged := &loggedResponseWriter{0, "", w}

		if IsSafeMethod(req.Method) || CheckReferer(req.Referer()) {
			mux.ServeHTTP(wLogged, req)
		} else {
			http.Error(wLogged, "CSRF failure", http.StatusForbidden)
		}

		logRequest(wLogged, req)
	})
	return newMux
}
