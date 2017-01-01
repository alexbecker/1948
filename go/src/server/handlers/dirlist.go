package handlers

import (
	"server/middleware"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func DirListHandler(baseFilepath http.Dir) http.Handler {
	return middleware.LoggedHandler(func(w http.ResponseWriter, req *http.Request) error {
		contents := make([]string, 0, 10)
		err := filepath.Walk(string(baseFilepath), filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != string(baseFilepath) {
				contents = append(contents, path)
			}
			return nil
		}))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(contents)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(encoded)
		return nil
	})
}
