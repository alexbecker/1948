package handlers

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var compressible = map[string]bool{
	".html": true,
	".css":  true,
	".js":   true,
}
var encodings = []string{"br", "gzip"}
var extensions = map[string]string{
	"br":   ".br",
	"gzip": ".gz",
}

func ServeEncFile(w http.ResponseWriter, req *http.Request, filename string) {
	if compressible[filepath.Ext(filename)] {
		for _, encoding := range encodings {
			accepts := strings.Split(req.Header.Get("Accept-Encoding"), ", ")
			for _, accept := range accepts {
				if encoding == accept {
					w.Header().Set("Content-Encoding", encoding)
					w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
					ServeFile(w, req, filename+extensions[encoding])
					return
				}
			}
		}
	}

	ServeFile(w, req, filename)
}

func SingleEncFileHandler(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ServeEncFile(w, req, filename)
	})
}

func ServeEncFiles(w http.ResponseWriter, req *http.Request, baseFilepath string) {
	path := filepath.Join(baseFilepath, req.URL.Path)
	ServeEncFile(w, req, path)
}

func EncHandler(baseFilepath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ServeEncFiles(w, req, baseFilepath)
	})
}
