package session

import (
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"github.com/albakov/go-weather-viewer/internal/storage/session"
	"github.com/albakov/go-weather-viewer/internal/storage/user"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	c := config.MustNew(filepath.Join(basePath, "/../../../config/app.toml"))
	c.DSN = c.DSNTest
	db := storage.MustNew(c)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	userStorage := user.NewStorage(db)
	sessionStorage := session.NewStorage(db)

	mockUser, _ := userStorage.Create("demo", "demo")

	mockSession := entity.Session{
		Id:        "f228b56e-b9e5-4e36-ae6e-75a047eecdc2",
		UserId:    mockUser.Id,
		ExpiresAt: "2030-01-01T00:00:00",
	}

	sessionStorage.Create(mockSession)

	s := NewService(c, db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:     Key,
		Value:    mockSession.Id,
		Path:     "/",
		Expires:  mockSession.ExpiresAtTime(),
		HttpOnly: false,
	})

	authenticated := s.IsAuthenticated(r)

	userStorage.Delete(mockUser.Id)
	sessionStorage.Delete(mockSession.Id)

	if !authenticated {
		t.Errorf("IsAuthenticated returned false")
	}
}

func TestAuthenticatedUser(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	c := config.MustNew(filepath.Join(basePath, "/../../../config/app.toml"))
	c.DSN = c.DSNTest
	db := storage.MustNew(c)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	userStorage := user.NewStorage(db)
	sessionStorage := session.NewStorage(db)

	mockUser, _ := userStorage.Create("demo", "demo")

	mockSession := entity.Session{
		Id:        "f228b56e-b9e5-4e36-ae6e-75a047eecdc2",
		UserId:    mockUser.Id,
		ExpiresAt: "2030-01-01T00:00:00",
	}

	sessionStorage.Create(mockSession)

	s := NewService(c, db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:     Key,
		Value:    mockSession.Id,
		Path:     "/",
		Expires:  mockSession.ExpiresAtTime(),
		HttpOnly: false,
	})

	authenticatedUser := s.AuthenticatedUser(r)

	userStorage.Delete(mockUser.Id)
	sessionStorage.Delete(mockSession.Id)

	if mockUser.Id != authenticatedUser.Id {
		t.Errorf("TestAuthenticatedUser: expected id %d, got %d", mockUser.Id, authenticatedUser.Id)
	}

	if mockUser.Login != authenticatedUser.Login {
		t.Errorf("TestAuthenticatedUser: expected login %s, got %s", mockUser.Login, authenticatedUser.Login)
	}
}

func TestAuthenticatedUserExpired(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	c := config.MustNew(filepath.Join(basePath, "/../../../config/app.toml"))
	c.DSN = c.DSNTest
	db := storage.MustNew(c)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	userStorage := user.NewStorage(db)
	sessionStorage := session.NewStorage(db)

	mockUser, _ := userStorage.Create("demo", "demo")

	mockSession := entity.Session{
		Id:        "f228b56e-b9e5-4e36-ae6e-75a047eecdc2",
		UserId:    mockUser.Id,
		ExpiresAt: "2022-01-01T00:00:00",
	}

	sessionStorage.Create(mockSession)

	s := NewService(c, db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:     Key,
		Value:    mockSession.Id,
		Path:     "/",
		Expires:  mockSession.ExpiresAtTime(),
		HttpOnly: false,
	})

	authenticatedUser := s.AuthenticatedUser(r)

	userStorage.Delete(mockUser.Id)
	sessionStorage.Delete(mockSession.Id)

	if authenticatedUser.Id != 0 {
		t.Errorf("TestAuthenticatedUserExpired: expected id %d, got %d", 0, authenticatedUser.Id)
	}
}
