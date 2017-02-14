package handlers

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"server/middleware"
	"strings"
)

var compressible = map[string]bool{
	".html": true,
	".css":  true,
	".js":   true,
}
var extensions = map[string]string{
	"br":   ".br",
	"gzip": ".gz",
}

// Attempts to serve filename, setting encoding and mimetype on success if passed.
// If filename is not found, returns false to allow the caller to handle appropriately.
func serveFile(w http.ResponseWriter, req *http.Request, filename, encoding, mimeType string) bool {
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
	if stats.IsDir() {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return true
	}

	file, err := os.Open(filename)
	if err != nil {
		middleware.Handle500(w, err)
		return true
	}
	defer file.Close()

	if encoding != "" {
		w.Header().Set("Content-Encoding", encoding)
		w.Header().Set("Content-Type", mimeType)
	}
	http.ServeContent(w, req, filename, stats.ModTime(), file)
	return true
}

// Serves an unencoded file.
func ServeFile(w http.ResponseWriter, req *http.Request, filename string) {
	found := serveFile(w, req, filename, "", "")
	if !found {
		http.NotFound(w, req)
	}
}

func SingleFileHandler(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ServeFile(w, req, filename)
	})
}

// Attempts to serve an encoded version of filename, if present and accepted.
func ServeEncFile(w http.ResponseWriter, req *http.Request, filename string) {
	if compressible[filepath.Ext(filename)] {
		for encoding, extension := range extensions {
			accepts := strings.Split(req.Header.Get("Accept-Encoding"), ", ")
			for _, accept := range accepts {
				if encoding == accept {
					mimeType := mime.TypeByExtension(filepath.Ext(filename))
					found := serveFile(w, req, filename+extension, encoding, mimeType)
					if found {
						return
					}
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
