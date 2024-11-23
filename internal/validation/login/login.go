package login

import (
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/validation"
	"net/http"
)

type Request struct {
	r            *http.Request
	fields       []entity.FormField
	errorMessage string
}

func NewValidator(r *http.Request) *Request {
	return &Request{
		r: r,
		fields: []entity.FormField{
			{
				Label: "Login",
				Name:  "login",
			},
			{
				Label: "Password",
				Name:  "password",
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

		cc.fields[i].Value = v
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
