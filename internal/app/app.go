package app

import (
	"database/sql"
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/controller/index"
	"github.com/albakov/go-weather-viewer/internal/controller/locations"
	"github.com/albakov/go-weather-viewer/internal/controller/login"
	"github.com/albakov/go-weather-viewer/internal/controller/register"
	"github.com/albakov/go-weather-viewer/internal/middleware"
	"github.com/albakov/go-weather-viewer/internal/services/session"
	"net/http"
	"strings"
)

type App struct {
	mux                *http.ServeMux
	config             *config.Config
	indexController    *index.Controller
	locationController *locations.Controller
	loginController    *login.Controller
	registerController *register.Controller
	sessionService     *session.Service
}

func New(config *config.Config, db *sql.DB) *App {
	commonController := controller.New()

	return &App{
		mux:                http.NewServeMux(),
		config:             config,
		indexController:    index.New(commonController, config, db),
		locationController: locations.New(commonController, config, db),
		loginController:    login.New(commonController, config, db),
		registerController: register.New(commonController, config, db),
		sessionService:     session.NewService(config, db),
	}
}

func (a *App) MustStart() {
	a.setRoutes()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.config.Host, a.config.Port), a)
	if err != nil {
		panic(err)
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := middleware.SetAuthenticatedUser(a.mux, a.sessionService, []string{"/css", "/js", "/images"})
	handler.ServeHTTP(w, r)
}

func (a *App) setRoutes() {
	a.setAssets()
	a.mux.HandleFunc("/", a.indexController.IndexHandler)
	a.mux.HandleFunc("/locations", a.locationController.IndexHandler)
	a.mux.HandleFunc("/locations/{id}/delete", a.locationController.DeleteHandler)
	a.mux.Handle("/login", middleware.NotAuthenticated(http.HandlerFunc(a.loginController.LoginHandler), a.sessionService))
	a.mux.Handle("/logout", middleware.Authenticated(http.HandlerFunc(a.loginController.LogoutHandler), a.sessionService))
	a.mux.Handle("/registration", middleware.NotAuthenticated(http.HandlerFunc(a.registerController.RegisterHandler), a.sessionService))
}

func (a *App) setAssets() {
	a.mux.Handle("/css/", a.noDirListing(http.FileServer(http.Dir("./view"))))
	a.mux.Handle("/js/", a.noDirListing(http.FileServer(http.Dir("./view"))))
	a.mux.Handle("/images/", a.noDirListing(http.FileServer(http.Dir("./view"))))
}

func (a *App) noDirListing(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/404", http.StatusFound)

			return
		}

		h.ServeHTTP(w, r)
	}
}
