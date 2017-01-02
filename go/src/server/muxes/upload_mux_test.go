package muxes

import (
	"github.com/DATA-DOG/go-sqlmock"
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
	mockDB     sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
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

	// Set up mock database for authentication.
	mockDB = auth.InitMock("salt", "secret")
	authCookie = auth.AuthCookie("user")

	result := m.Run()
	os.Exit(result)
}

func expectRoleLookup() {
	mockDB.ExpectQuery("SELECT rowid FROM users WHERE username = ?").
		WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"rowid"}).AddRow(1))
	mockDB.ExpectQuery("SELECT role FROM user_roles WHERE userid = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("upload"))
	mockDB.ExpectQuery("SELECT childid FROM user_inheritance WHERE parentid = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"childid"}))
}

func checkSqlExpectations(t *testing.T) {
	err := mockDB.ExpectationsWereMet()
	if err != nil {
		t.Errorf("sql error: %s", err)
	}
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
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusOK, "test", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
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
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusOK, `["/test","/tmp.html"]`, response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
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
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusCreated, "File /test created", response)

	request = httptest.NewRequest("GET", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusOK, "content2", response)

	request = httptest.NewRequest("DELETE", "/uploads/test", nil)
	request.AddCookie(&authCookie)
	response = httptest.NewRecorder()
	expectRoleLookup()
	mux.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusNoContent, "", response)
}
