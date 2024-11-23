package user

import (
	"errors"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"log"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCreateUserWithNotUniqueLogin(t *testing.T) {
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

	userStorage := NewStorage(db)

	login := "demo"
	password := "demo"

	mockUser, _ := userStorage.Create(login, password)
	_, err := userStorage.Create(login, password)

	userStorage.Delete(mockUser.Id)

	if err == nil {
		t.Errorf("Create user with not unique login should fail")
	}

	if !errors.Is(err, storage.EntityAlreadyExistsError) {
		t.Errorf("Create user with not unique login should fail")
	}
}
