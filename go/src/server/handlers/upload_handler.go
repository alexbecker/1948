package handlers

import (
	"server/middleware"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(baseFilepath http.Dir, methods map[string]bool, overwrite bool) http.Handler {
	return middleware.LoggedHandler(func(w http.ResponseWriter, req *http.Request) error {
		if !methods[req.Method] {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return nil
		}

		switch req.Method {
		case "GET":
			path := filepath.Join(string(baseFilepath), req.URL.Path)
			ServeFile(w, req, path)
		case "POST":
			uploadFile, uploadFileHeader, uploadErr := req.FormFile("file")
			if uploadErr != nil {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return nil
			}

			uploadFilename := filepath.Join(req.URL.Path, uploadFileHeader.Filename)

			if !overwrite {
				// check that file does not already exist
				_, err := baseFilepath.Open(uploadFilename)
				if err == nil {
					http.Error(w, "409 conflict", http.StatusConflict)
					return nil
				} else if !os.IsNotExist(err) {
					return err
				}
			}

			fp, creationErr := os.Create(filepath.Join(string(baseFilepath), uploadFilename))
			if creationErr != nil {
				return nil
			}

			defer fp.Close()

			_, copyErr := io.Copy(fp, uploadFile)
			if copyErr != nil {
				return nil
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("File " + uploadFilename + " created"))
		case "DELETE":
			path := filepath.Join(string(baseFilepath), req.URL.Path)
			err := os.Remove(path)
			if os.IsNotExist(err) {
				http.NotFound(w, req)
				return nil
			} else if err != nil {
				return err
			}

			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("File " + path + " deleted"))
		}

		return nil
	})
}
