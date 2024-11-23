package login

import (
	"net/http"
	"net/url"
	"testing"
)

func getRequestForTest(form url.Values) *http.Request {
	return &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/login",
		},
		Host: "localhost",
		Form: form,
	}
}

func TestEmptyFields(t *testing.T) {
	form := url.Values{}
	form.Add("login", "")
	form.Add("password", "")

	validator := NewValidator(getRequestForTest(form))
	validator.Validate()

	if validator.IsValid() {
		t.Errorf("IsValid() should return false")
	}
}

func TestFieldsOk(t *testing.T) {
	form := url.Values{}
	form.Add("login", "john")
	form.Add("password", "doe")

	validator := NewValidator(getRequestForTest(form))
	validator.Validate()

	if !validator.IsValid() {
		t.Errorf("IsValid() should return true")
	}
}
