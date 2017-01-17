package middleware

import (
	"log"
	"net/http"
	"os"
)

var InfoLogger, ErrLogger *log.Logger

func init() {
	InfoLogger = log.New(os.Stdin, "INFO: ", log.LstdFlags|log.LUTC)
	ErrLogger = log.New(os.Stderr, "ERR: ", log.LstdFlags|log.LUTC)
}

type loggedResponseWriter struct {
	status         int
	strippedPrefix string
	http.ResponseWriter
}

func (w *loggedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggedResponseWriter) Write(p []byte) (int, error) {
	// If WriteHeader hasn't been called, the status is implicitly set to 200 OK.
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(p)
}

type LoggedHandler func(http.ResponseWriter, *http.Request) error

func Handle500(w http.ResponseWriter, err error) {
	ErrLogger.Print(err)
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

func (handler LoggedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := handler(w, req)
	if err != nil {
		Handle500(w, err)
	}
}

func StripPrefix(prefix string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedW, isLogged := w.(*loggedResponseWriter)
		if isLogged {
			loggedW.strippedPrefix = prefix
		}
		http.StripPrefix(prefix, handler).ServeHTTP(w, r)
	})
}

func logRequest(w *loggedResponseWriter, req *http.Request) {
	referer := req.Referer()
	if referer == "" {
		referer = "-"
	}
	userAgent := req.UserAgent()
	if userAgent == "" {
		userAgent = "-"
	}

	InfoLogger.Printf("%d %s %s%s %s %s", w.status, req.Method, w.strippedPrefix, req.URL.String(), referer, userAgent)
}
