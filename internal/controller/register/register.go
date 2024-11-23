package register

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/services/password"
	"github.com/albakov/go-weather-viewer/internal/services/session"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"github.com/albakov/go-weather-viewer/internal/storage/user"
	"github.com/albakov/go-weather-viewer/internal/validation/register"
	"net/http"
)

type Controller struct {
	commonController *controller.Controller
	config           *config.Config
	userStorage      *user.Storage
	sessionService   *session.Service
}

func New(commonController *controller.Controller, config *config.Config, db *sql.DB) *Controller {
	return &Controller{
		commonController: commonController,
		config:           config,
		userStorage:      user.NewStorage(db),
		sessionService:   session.NewService(config, db),
	}
}

func (cc *Controller) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.commonController.ShowResponse(
			w,
			"register.html",
			entity.PageData{
				PageTitle: "Registration",
				Header:    cc.commonController.GetHeaderData(r),
				Data:      cc.formRegistration(),
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

func (cc *Controller) authHandler(w http.ResponseWriter, r *http.Request) {
	validator := register.NewValidator(r, cc.config)
	validator.Validate()

	if !validator.IsValid() {
		cc.registerFailed(w, validator.ErrorMessage())

		return
	}

	hashedPassword, err := password.CreateHashedPassword(validator.ValueByName("password"))
	if err != nil {
		cc.registerFailed(w, controller.MessageServerError)

		return
	}

	userItem, err := cc.userStorage.Create(validator.ValueByName("login"), hashedPassword)
	if err != nil {
		if errors.Is(err, storage.EntityAlreadyExistsError) {
			cc.registerFailed(w, controller.MessageLoginAlreadyExists)

			return
		}

		cc.registerFailed(w, controller.MessageServerError)

		return
	}

	sessionItem, err := cc.sessionService.Create(userItem.Id)
	if err != nil {
		cc.registerFailed(w, controller.MessageServerError)

		return
	}

	cc.commonController.SetCookie(w, session.Key, sessionItem.Id, sessionItem.ExpiresAtTime())
	cc.commonController.RedirectTo(w, "/")
}

func (cc *Controller) registerFailed(w http.ResponseWriter, message string) {
	fr := cc.formRegistration()
	fr.Message = message

	cc.commonController.BackWithError(
		w,
		"register.html",
		entity.PageData{
			PageTitle: "Registration",
			Header:    cc.commonController.GetHeaderData(&http.Request{}),
			Data:      fr,
		},
	)
}

func (cc *Controller) formRegistration() entity.FormRegistration {
	return entity.FormRegistration{
		LoginPlaceholder:    fmt.Sprintf("Login, minimum %d characters", cc.config.FieldValueMinLength),
		PasswordPlaceholder: fmt.Sprintf("Password, minimum %d characters", cc.config.FieldValueMinLength),
	}
}
