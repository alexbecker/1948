package testingutils

import (
	"net/http/httptest"
	"testing"
)

func ExpectSuccess(t *testing.T, code int, content string, response *httptest.ResponseRecorder) {
	if response.Code != code {
		t.Errorf("Response code %d != %d", response.Code, code)
	}
	responseString := response.Body.String()
	if responseString != content {
		t.Errorf("Response %s != %s", responseString, content)
	}
}

func ExpectFailure(t *testing.T, code int, response *httptest.ResponseRecorder) {
	if response.Code != code {
		t.Errorf("Response code %d != %d", response.Code, code)
	}
}
