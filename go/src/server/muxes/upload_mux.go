package muxes

import (
	"net/http"
	"server/auth"
	"server/handlers"
	"server/middleware"
)

func HandleUploadPage(mux *http.ServeMux, path, role, page string, baseDir http.Dir, overwrite bool) {
	dirMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"DELETE": true,
	}
	dirHandler := handlers.UploadHandler(baseDir, dirMethods, overwrite)
	dirHandler = auth.AuthHandler(role, dirHandler)
	dirHandler = middleware.StripPrefix(path, dirHandler)
	mux.Handle(path+"/", dirHandler)

	dirListHandler := handlers.DirListHandler(baseDir)
	dirListHandler = auth.AuthHandler(role, dirListHandler)
	mux.Handle(path, dirListHandler)

	pageHandler := handlers.SingleEncFileHandler(page)
	pageHandler = auth.AuthHandler(role, pageHandler)
	mux.Handle(path+".html", pageHandler)
}
