package handlers

import (
	"net/http"
	"os"
	"server/middleware"
)

// Attempt to serve filename, returning false if it does not exist.
func MaybeServeFile(w http.ResponseWriter, req *http.Request, filename string) bool {
	stats, err := os.Stat(filename)
	if err != nil {
		if os.IsPermission(err) {
			http.Error(w, "403 forbidden", http.StatusForbidden)
		} else if os.IsNotExist(err) {
			return false
		} else {
			middleware.Handle500(w, err)
		}
		return true
	}
	file, err := os.Open(filename)
	if err != nil {
		middleware.Handle500(w, err)
		return true
	}
	defer file.Close()

	http.ServeContent(w, req, filename, stats.ModTime(), file)
	return true
}

func ServeFile(w http.ResponseWriter, req *http.Request, filename string) {
	ok := MaybeServeFile(w, req, filename)
	if !ok {
		http.NotFound(w, req)
	}
}

func SingleFileHandler(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ServeFile(w, req, filename)
	})
}
