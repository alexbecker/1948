package handlers

import (
	"net/http/httptest"
	"testing"
)

func TestListing(t *testing.T) {
	handler := DirListHandler("example_dir")
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Errorf("Response code %d != 200", response.Code)
	}
	responseString := response.Body.String()
	expectedResponseString := "[\"example_dir/a\",\"example_dir/b\",\"example_dir/b/c\"]"
	if responseString != expectedResponseString {
		t.Errorf("Response %s != %s", responseString, expectedResponseString)
	}
}
