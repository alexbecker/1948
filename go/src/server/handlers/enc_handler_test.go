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

var handler http.Handler

func TestMain(m *testing.M) {
	tempDir, err := ioutil.TempDir("", "enc_handler_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir) // clean up
	tempFile := filepath.Join(tempDir, "tmp.html")
	err = ioutil.WriteFile(tempFile, []byte("test"), 0666)
	if err != nil {
		panic(err)
	}
	handler = EncHandler(tempDir)
	os.Exit(m.Run())
}

func TestServeEncFilesNoAccept(t *testing.T) {
	request := httptest.NewRequest("GET", "/tmp.html", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectSuccess(t, 200, "test", response)
}

func TestServeEncFilesAcceptFallback(t *testing.T) {
	request := httptest.NewRequest("GET", "/tmp.html", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectSuccess(t, 200, "test", response)
}
