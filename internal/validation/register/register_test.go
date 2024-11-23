package register

import (
	"github.com/albakov/go-weather-viewer/internal/config"
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
	form.Add("password_again", "")

	validator := NewValidator(getRequestForTest(form), &config.Config{FieldValueMinLength: 6})
	validator.Validate()

	if validator.IsValid() {
		t.Errorf("IsValid() should return false")
	}
}

func TestFieldsOk(t *testing.T) {
	form := url.Values{}
	form.Add("login", "john_doe")
	form.Add("password", "john_doe_123")
	form.Add("password_again", "john_doe_123")

	validator := NewValidator(getRequestForTest(form), &config.Config{FieldValueMinLength: 6})
	validator.Validate()

	if !validator.IsValid() {
		t.Errorf("IsValid() should return true")
	}
}
