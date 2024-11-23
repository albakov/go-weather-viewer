package controller

import (
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/util"
	"html/template"
	"net/http"
	"time"
)

const f = "controller.Controller"

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) ShowResponse(w http.ResponseWriter, templatePath string, data entity.PageData) {
	const op = "ShowResponse"

	c.setHeaders(w)
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFiles(
		fmt.Sprintf("view/%s", templatePath),
		"view/header.html",
		"view/search-block.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, data)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowNotFound(w http.ResponseWriter, r *http.Request) {
	const op = "ShowNotFound"

	c.setHeaders(w)
	w.WriteHeader(http.StatusNotFound)

	tmpl := template.Must(template.ParseFiles(
		"view/error.html",
		"view/header.html",
		"view/search-block.html",
		"view/footer.html",
	))
	err := tmpl.Execute(
		w,
		entity.PageData{
			PageTitle: "Not Found",
			Header:    c.GetHeaderData(r),
			Data:      entity.ErrorData{Code: http.StatusNotFound},
		},
	)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) BackWithError(w http.ResponseWriter, templatePath string, data entity.PageData) {
	const op = "BackWithError"

	c.setHeaders(w)
	w.WriteHeader(http.StatusUnprocessableEntity)

	tmpl := template.Must(template.ParseFiles(
		fmt.Sprintf("view/%s", templatePath),
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, data)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	const op = "ShowMethodNotAllowedError"

	c.setHeaders(w)
	w.WriteHeader(http.StatusMethodNotAllowed)

	tmpl := template.Must(template.ParseFiles(
		"view/error.html",
		"view/header.html",
		"view/search-block.html",
		"view/footer.html",
	))
	err := tmpl.Execute(
		w,
		entity.PageData{
			PageTitle: MessageMethodNotAllowed,
			Header:    c.GetHeaderData(r),
			Data:      entity.ErrorData{Code: http.StatusMethodNotAllowed},
		},
	)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) RedirectTo(w http.ResponseWriter, url string) {
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (c *Controller) SetCookie(w http.ResponseWriter, name, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})
}

func (c *Controller) setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
}

func (c *Controller) GetHeaderData(r *http.Request) entity.Header {
	value := r.Context().Value("user")
	if value == nil {
		value = entity.User{}
	}

	return entity.Header{User: value.(entity.User)}
}
