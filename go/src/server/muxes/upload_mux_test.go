package muxes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server/auth"
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

func expectSuccess(t *testing.T, code int, content string, response *httptest.ResponseRecorder) {
	if response.Code != code {
		t.Errorf("Response code %d != %d", response.Code, code)
	}
	responseString := response.Body.String()
	if responseString != content {
		t.Errorf("Response %s != %s", responseString, content)
	}
}

func expectFailure(t *testing.T, code int, response *httptest.ResponseRecorder) {
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Response code %d != http.StatusUnauthorized", response.Code)
	}
}

func TestAuthFailures(t *testing.T) {
	mux := http.NewServeMux()
	HandleUploadPage(mux, "/uploads", "upload", uploadPage, uploadDir, true)

	request := httptest.NewRequest("GET", "/uploads.html", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectFailure(t, http.StatusUnauthorized, response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectFailure(t, http.StatusUnauthorized, response)

	request = httptest.NewRequest("GET", "/uploads/a", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectFailure(t, http.StatusUnauthorized, response)

	request = httptest.NewRequest("POST", "/uploads/", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectFailure(t, http.StatusUnauthorized, response)

	request = httptest.NewRequest("DELETE", "/uploads/a", nil)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectFailure(t, http.StatusUnauthorized, response)
}

func TestHappyPath(t *testing.T) {
	mux := http.NewServeMux()
	HandleUploadPage(mux, "/uploads", "upload", uploadPage, uploadDir, true)

	request := httptest.NewRequest("GET", "/uploads.html", nil)
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusOK, "test", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusOK, `["/tmp.html"]`, response)

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
	expectSuccess(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusOK, `["/test","/tmp.html"]`, response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusOK, "content", response)

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
	expectSuccess(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusOK, "content2", response)

	request = httptest.NewRequest("DELETE", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	expectSuccess(t, http.StatusNoContent, "", response)
}
