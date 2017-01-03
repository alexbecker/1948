package auth

import (
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"server/testingutils"
	"testing"
)

var (
	handler http.Handler
	mockDB  sqlmock.Sqlmock
)

func init() {
	mockDB = InitMock("salt", "secret")
	handler = AuthHandler("role", http.NotFoundHandler())
}

func expectPasswordLookup(username string) {
	rows := sqlmock.NewRows([]string{"hash"})
	if username == "user" {
		rows.AddRow(hash(username, "password"))
	}
	mockDB.ExpectQuery("SELECT hash FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(rows)
}

func expectRoleLookup(username string) {
	roles := sqlmock.NewRows([]string{"role"})
	if username == "user" {
		roles.AddRow("role")
	}
	mockDB.ExpectQuery("SELECT rowid FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"rowid"}).AddRow(1))
	mockDB.ExpectQuery("SELECT role FROM user_roles WHERE userid = ?").
		WithArgs(1).
		WillReturnRows(roles)
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

func TestChallengeIssued(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
	challenge := response.Header().Get("WWW-Authenticate")
	expectedChallenge := `Basic realm="role"`
	if challenge != expectedChallenge {
		t.Errorf("Authentication challenge %s != %s", challenge, expectedChallenge)
	}
}

func TestChallengePassed(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.SetBasicAuth("user", "password")
	response := httptest.NewRecorder()
	expectPasswordLookup("user")
	handler.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusNotFound, "404 page not found\n", response)
	cookie := *response.Result().Cookies()[0]
	expectedCookie := AuthCookie("user")
	if cookie.Name != expectedCookie.Name || cookie.Value != expectedCookie.Value {
		t.Errorf("Auth cookie %s != %s", cookie.String(), expectedCookie.String())
	}
}

func TestChallengeFailed(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.SetBasicAuth("user", "password2")
	response := httptest.NewRecorder()
	expectPasswordLookup("user")
	handler.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
	if len(response.Result().Cookies()) > 0 {
		t.Errorf("Unexpected cookie set")
	}
}

func TestMissingUser(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.SetBasicAuth("user2", "password")
	response := httptest.NewRecorder()
	expectPasswordLookup("user2")
	handler.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
	if len(response.Result().Cookies()) > 0 {
		t.Errorf("Unexpected cookie set")
	}
}

func TestAuthCookieAccepted(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	authCookie := AuthCookie("user")
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	expectRoleLookup("user")
	handler.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusNotFound, "404 page not found\n", response)
}

func TestAuthCookieUnsigned(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	authCookie := http.Cookie{
		Name:  "auth",
		Value: "user",
	}
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
}

func TestAuthCookieBogusSignature(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	authCookie := http.Cookie{
		Name:  "auth",
		Value: "user.aaaaaaaaaaaaaaaaaaa",
	}
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	testingutils.ExpectResponse(t, http.StatusUnauthorized, "", response)
}

func TestAuthCookieWrongRole(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	authCookie := AuthCookie("user2")
	request.AddCookie(&authCookie)
	response := httptest.NewRecorder()
	expectRoleLookup("user2")
	handler.ServeHTTP(response, request)
	checkSqlExpectations(t)
	testingutils.ExpectResponse(t, http.StatusForbidden, "403 forbidden\n", response)
}
