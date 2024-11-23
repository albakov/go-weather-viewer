package session

import (
	"database/sql"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/storage/session"
	"github.com/albakov/go-weather-viewer/internal/storage/user"
	"github.com/albakov/go-weather-viewer/internal/util"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	f   = "session"
	Key = "sessionId"
)

type Service struct {
	config         *config.Config
	sessionStorage *session.Storage
	userStorage    *user.Storage
}

func NewService(config *config.Config, db *sql.DB) *Service {
	return &Service{
		config:         config,
		sessionStorage: session.NewStorage(db),
		userStorage:    user.NewStorage(db),
	}
}

func (s *Service) Create(userId int64) (entity.Session, error) {
	const op = "Create"

	sessionItem := entity.Session{
		Id:        uuid.NewString(),
		UserId:    userId,
		ExpiresAt: time.Now().UTC().Add(time.Duration(s.config.SessionPeriod) * time.Hour).Format(time.DateTime),
	}

	err := s.sessionStorage.Create(sessionItem)
	if err != nil {
		util.LogError(f, op, err)

		return entity.Session{}, err
	}

	return sessionItem, nil
}

func (s *Service) Delete(r *http.Request) {
	const op = "Delete"

	cookie, err := r.Cookie(Key)
	if err != nil || cookie.Value == "" {
		return
	}

	err = s.sessionStorage.Delete(cookie.Value)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (s *Service) AuthenticatedUser(r *http.Request) entity.User {
	const op = "AuthenticatedUser"

	cookie, err := r.Cookie(Key)
	if err != nil || cookie.Value == "" {
		return entity.User{}
	}

	sessionItem, err := s.sessionStorage.GetById(cookie.Value)
	if err != nil {
		return entity.User{}
	}

	if sessionItem.ExpiresAtTime().Before(time.Now().UTC()) {
		err := s.sessionStorage.Delete(sessionItem.Id)
		if err != nil {
			util.LogError(f, op, err)
		}

		return entity.User{}
	}

	userItem, err := s.userStorage.GetById(sessionItem.UserId)
	if err != nil {
		return entity.User{}
	}

	return userItem
}

func (s *Service) IsAuthenticated(r *http.Request) bool {
	const op = "IsAuthenticated"

	cookie, err := r.Cookie(Key)
	if err != nil || cookie.Value == "" {
		return false
	}

	sessionItem, err := s.sessionStorage.GetById(cookie.Value)
	if err != nil {
		return false
	}

	if sessionItem.ExpiresAtTime().Before(time.Now().UTC()) {
		err := s.sessionStorage.Delete(sessionItem.Id)
		if err != nil {
			util.LogError(f, op, err)
		}

		return false
	}

	return s.userStorage.IsExistsById(sessionItem.UserId)
}
