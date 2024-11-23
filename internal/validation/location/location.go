package location

import (
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/validation"
	"net/http"
	"strconv"
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
				Label: "name",
				Name:  "name",
				Value: "",
			},
			{
				Label: "latitude",
				Name:  "latitude",
				Value: "",
			},
			{
				Label: "longitude",
				Name:  "longitude",
				Value: "",
			},
		},
	}
}

func (r *Request) Validate() {
	for i, field := range r.fields {
		v := r.r.FormValue(field.Name)
		if v == "" {
			r.errorMessage = fmt.Sprintf(validation.MessageFieldEmpty, field.Label)

			return
		}

		r.fields[i].Value = v
	}
}

func (r *Request) IsValid() bool {
	return r.errorMessage == ""
}

func (r *Request) ErrorMessage() string {
	return r.errorMessage
}

func (r *Request) ValueByName(name string) string {
	for _, field := range r.fields {
		if field.Name == name {
			return field.Value
		}
	}

	return ""
}

func (r *Request) GetFloat(key string) float64 {
	v, err := strconv.ParseFloat(r.ValueByName(key), 64)
	if err != nil {
		return 0.0
	}

	return v
}
