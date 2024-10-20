package handlers

import (
	"compress/gzip"
	"github.com/alexbecker/1948/testingutils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var handler http.Handler

func TestMain(m *testing.M) {
	tempDir, _ := ioutil.TempDir("", "enc_handler_test")
	defer os.RemoveAll(tempDir) // clean up
	tempFile := filepath.Join(tempDir, "tmp.html")
	ioutil.WriteFile(tempFile, []byte("test"), 0666)
	fp, _ := os.Create(filepath.Join(tempDir, "tmp2.html.gz"))
	gzipWriter := gzip.NewWriter(fp)
	gzipWriter.Write([]byte("test2"))
	gzipWriter.Close()
	fp.Close()
	handler = EncHandler(tempDir)
	os.Exit(m.Run())
}

func TestServeEncFilesNoAccept(t *testing.T) {
	request := httptest.NewRequest("GET", "/tmp.html", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, "test", response)
}

func TestServeEncFilesAcceptFallback(t *testing.T) {
	request := httptest.NewRequest("GET", "/tmp.html", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, "test", response)
}

func TestServeEncFilesAcceptGzip(t *testing.T) {
	request := httptest.NewRequest("GET", "/tmp2.html", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("Response code %d != %d", response.Code, http.StatusOK)
	}
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		t.Errorf("Error unzipping response: %s", err)
	}
	content, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		t.Errorf("Error unzipping response: %s", err)
	}
	if string(content) != "test2" {
		t.Errorf("Response %s != test2", content)
	}
}

func TestNotFound(t *testing.T) {
	request := httptest.NewRequest("GET", "/404", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusNotFound, "404 page not found\n", response)
}

func TestNotServeDirectory(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusNotFound, "404 page not found\n", response)
}
