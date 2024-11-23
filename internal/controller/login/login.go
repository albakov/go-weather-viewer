package login

import (
	"database/sql"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/services/password"
	"github.com/albakov/go-weather-viewer/internal/services/session"
	"github.com/albakov/go-weather-viewer/internal/storage/user"
	"github.com/albakov/go-weather-viewer/internal/validation/login"
	"net/http"
	"time"
)

type Controller struct {
	commonController *controller.Controller
	userStorage      *user.Storage
	sessionService   *session.Service
}

func New(commonController *controller.Controller, config *config.Config, db *sql.DB) *Controller {
	return &Controller{
		commonController: commonController,
		userStorage:      user.NewStorage(db),
		sessionService:   session.NewService(config, db),
	}
}

func (cc *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.commonController.ShowResponse(
			w,
			"login.html",
			entity.PageData{
				PageTitle: "Log in",
				Header:    cc.commonController.GetHeaderData(r),
				Data:      cc.formLogin(),
			},
		)

		return
	}

	if r.Method == http.MethodPost {
		cc.authHandler(w, r)

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w, r)
}

func (cc *Controller) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		cc.commonController.ShowMethodNotAllowedError(w, r)

		return
	}

	cc.sessionService.Delete(r)
	cc.commonController.SetCookie(w, session.Key, "", time.Time{})
	cc.commonController.RedirectTo(w, "/")
}

func (cc *Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	validator := login.NewValidator(r)
	validator.Validate()

	if !validator.IsValid() {
		cc.loginFailed(w)

		return
	}

	userItem, err := cc.userStorage.GetByLogin(validator.ValueByName("login"))
	if err != nil {
		cc.loginFailed(w)

		return
	}

	if !password.CheckPassword(validator.ValueByName("password"), userItem.Password) {
		cc.loginFailed(w)

		return
	}

	sessionItem, err := cc.sessionService.Create(userItem.Id)
	if err != nil {
		cc.loginFailed(w)

		return
	}

	cc.commonController.SetCookie(w, session.Key, sessionItem.Id, sessionItem.ExpiresAtTime())
	cc.commonController.RedirectTo(w, "/")
}

func (cc *Controller) loginFailed(w http.ResponseWriter) {
	fl := cc.formLogin()
	fl.Message = controller.MessageLoginOrPasswordNotValid

	cc.commonController.BackWithError(
		w,
		"login.html",
		entity.PageData{
			PageTitle: "Log in",
			Header:    cc.commonController.GetHeaderData(&http.Request{}),
			Data:      fl,
		},
	)
}

func (cc *Controller) formLogin() entity.FormLogin {
	return entity.FormLogin{
		LoginPlaceholder:    "Login",
		PasswordPlaceholder: "Password",
	}
}
