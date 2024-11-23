package register

import (
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/validation"
	"net/http"
)

type Request struct {
	r            *http.Request
	config       *config.Config
	fields       []entity.FormField
	errorMessage string
}

func NewValidator(r *http.Request, config *config.Config) *Request {
	return &Request{
		r:      r,
		config: config,
		fields: []entity.FormField{
			{
				Label: "Login",
				Name:  "login",
			},
			{
				Label: "Password",
				Name:  "password",
			},
			{
				Label: "Password Again",
				Name:  "password_again",
			},
		},
	}
}

func (cc *Request) Validate() {
	for i, field := range cc.fields {
		v := cc.r.FormValue(field.Name)
		if v == "" {
			cc.errorMessage = fmt.Sprintf(validation.MessageFieldEmpty, field.Label)

			return
		}

		if int64(len(v)) < cc.config.FieldValueMinLength {
			cc.errorMessage = fmt.Sprintf(validation.MessageFieldTooSmall, field.Label, cc.config.FieldValueMinLength)

			return
		}

		cc.fields[i].Value = v
	}

	if cc.ValueByName("password") != cc.ValueByName("password_again") {
		cc.errorMessage = validation.MessagePasswordRepeatInvalid

		return
	}
}

func (cc *Request) IsValid() bool {
	return cc.errorMessage == ""
}

func (cc *Request) ErrorMessage() string {
	return cc.errorMessage
}

func (cc *Request) ValueByName(name string) string {
	for _, field := range cc.fields {
		if field.Name == name {
			return field.Value
		}
	}

	return ""
}
