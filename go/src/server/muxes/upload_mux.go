package muxes

import (
	"net/http"
	"server/auth"
	"server/handlers"
	"server/middleware"
)

// An opinionated implementation of the various handlers necessary for a read/write file server, including:
//   Serves uploadPage at path.html, which should make the following requests via forms and AJAX.
//   GETs to path will return a JSON array of all files in baseDir, relative to baseDir.
//   POSTs to path/ will upload a file to baseDir with the given name, overwriting any existing file by that name if overwrite is true.
//   GETs to path/filename will download baseDir/filename.
//   DELETEs to path/filename will delete baseDir/filename.
// These handlers all require the same role for uploading, downloading and deleting.
// Note that this places no limits on upload file size or type.
// Only trusted users should be given authorization for this.
func HandleUploadPage(mux *http.ServeMux, path, role, uploadPage string, baseDir http.Dir, overwrite bool) {
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

	pageHandler := handlers.SingleEncFileHandler(uploadPage)
	pageHandler = auth.AuthHandler(role, pageHandler)
	mux.Handle(path+".html", pageHandler)
}
