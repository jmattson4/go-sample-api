package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/jmattson4/go-sample-api/app"
)

var app a.App

func init() {
	app.InitializeTesting()
}

func TestGetNewsByWebName(t *testing.T) {
	body := &httpReader{}
	req, _ := http.NewRequest("GET", "/api/globalnews", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

//this function checks the response code of expected against actual and if they dont match then
//	it throws a testing error
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

//executes the request against the http server
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}
