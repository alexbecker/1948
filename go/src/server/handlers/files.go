package handlers

import (
	"server/middleware"
	"net/http"
	"os"
)

func ServeFile(w http.ResponseWriter, req *http.Request, filename string) {
	stats, err := os.Stat(filename)
	if err != nil {
		if os.IsPermission(err) {
			http.Error(w, "403 forbidden", http.StatusForbidden)
		} else if os.IsNotExist(err) {
			http.NotFound(w, req)
		} else {
			middleware.Handle500(w, err)
		}
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		middleware.Handle500(w, err)
		return
	}
	defer file.Close()

	http.ServeContent(w, req, filename, stats.ModTime(), file)
}

func SingleFileHandler(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ServeFile(w, req, filename)
	})
}
