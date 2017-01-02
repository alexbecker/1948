package muxes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server/auth"
	"server/testingutils"
	"strings"
	"testing"
)

var (
	authCookie http.Cookie
	uploadDir  http.Dir
	uploadPage string
)

func TestMain(m *testing.M) {
	os.Setenv("SECRET", "longsecrethahaha")

	// Create temporary directory for uploads and uploadPage to serve.
	uploadDirName, err := ioutil.TempDir("", "example")
	if err != nil {
		panic(err)
	}
	uploadDir = http.Dir(uploadDirName)
	defer os.RemoveAll(uploadDirName) // clean up
	uploadPage = filepath.Join(uploadDirName, "tmp.html")
	err = ioutil.WriteFile(uploadPage, []byte("test"), 0666)
	if err != nil {
		panic(err)
	}

	// Set up mock authentication.
	auth.SetMockUser("user", map[string]bool{
		"upload": true,
	})
	authCookie = auth.AuthCookie("user")

	result := m.Run()
	os.Exit(result)
}

func TestAuthFailures(t *testing.T) {
	mux := http.NewServeMux()
	HandleUploadPage(mux, "/uploads", "upload", uploadPage, uploadDir, true)

	request := httptest.NewRequest("GET", "/uploads.html", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)

	request = httptest.NewRequest("GET", "/uploads/a", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)

	request = httptest.NewRequest("POST", "/uploads/", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)

	request = httptest.NewRequest("DELETE", "/uploads/a", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
}

func TestHappyPath(t *testing.T) {
	mux := http.NewServeMux()
	HandleUploadPage(mux, "/uploads", "upload", uploadPage, uploadDir, true)

	request := httptest.NewRequest("GET", "/uploads.html", nil)
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, "test", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, `["/tmp.html"]`, response)

	request = httptest.NewRequest("POST", "/uploads/", strings.NewReader(`Content-Type: multipart/form-data; boundary=---------------------------18452080271361591831786946236
Content-Length: 23

-----------------------------18452080271361591831786946236
Content-Disposition: form-data; name="file"; filename="test"
Content-Type: application/octet-stream

content
-----------------------------18452080271361591831786946236--`))
	request.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------18452080271361591831786946236")
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, `["/test","/tmp.html"]`, response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, "content", response)

	request = httptest.NewRequest("POST", "/uploads/", strings.NewReader(`Content-Type: multipart/form-data; boundary=---------------------------18452080271361591831786946236
Content-Length: 24

-----------------------------18452080271361591831786946236
Content-Disposition: form-data; name="file"; filename="test"
Content-Type: application/octet-stream

content2
-----------------------------18452080271361591831786946236--`))
	request.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------18452080271361591831786946236")
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusOK, "content2", response)

	request = httptest.NewRequest("DELETE", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusNoContent, "", response)
}
