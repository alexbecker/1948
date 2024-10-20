package testingutils

import (
	"net/http/httptest"
	"testing"
)

func ExpectResponse(t *testing.T, code int, content string, response *httptest.ResponseRecorder) {
	if response.Code != code {
		t.Errorf("Response code %d != %d", response.Code, code)
	}
	encoding := response.Header().Get("Content-Encoding")
	if encoding != "" {
		t.Errorf("Unexpected encoded response (%s)", encoding)
	}
	responseString := response.Body.String()
	if responseString != content {
		t.Errorf("Response %s != %s", responseString, content)
	}
}
