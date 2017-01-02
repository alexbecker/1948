package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server/testingutils"
	"testing"
)

func TestListing(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "dirlist_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir) // clean up
	_, err = os.Create(filepath.Join(tempDir, "a"))
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(filepath.Join(tempDir, "b"), os.ModeDir|os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, err = os.Create(filepath.Join(tempDir, "b", "c"))
	if err != nil {
		panic(err)
	}

	handler := DirListHandler(http.Dir(tempDir))
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, `["/a","/b","/b/c"]`, response)
}
