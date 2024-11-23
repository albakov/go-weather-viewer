package user

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"github.com/albakov/go-weather-viewer/internal/util"
	"github.com/go-sql-driver/mysql"
)

const f = "user.Storage"

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(login, password string) (entity.User, error) {
	const op = "Create"

	stmt, err := s.db.Prepare("INSERT INTO users (login, password) VALUES (?, ?)")
	if err != nil {
		util.LogError(f, op, err)

		return entity.User{}, err
	}
	defer stmt.Close()

	exec, err := stmt.Exec(login, password)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.User{}, storage.EntityAlreadyExistsError
		}

		util.LogError(f, op, err)

		return entity.User{}, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		util.LogError(f, op, err)

		return entity.User{}, err
	}

	return s.GetById(id)
}

func (s *Storage) GetByLogin(login string) (entity.User, error) {
	const op = "GetByLogin"

	row := s.db.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	if err := row.Err(); err != nil {
		util.LogError(f, op, err)

		return entity.User{}, err
	}

	user := entity.User{}

	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.User{}, err
	}

	return user, nil
}

func (s *Storage) GetById(id int64) (entity.User, error) {
	const op = "GetById"

	row := s.db.QueryRow("SELECT id, login, password FROM users WHERE id = ?", id)
	if err := row.Err(); err != nil {
		util.LogError(f, op, err)

		return entity.User{}, err
	}

	user := entity.User{}

	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.User{}, err
	}

	return user, nil
}

func (s *Storage) IsExistsById(id int64) bool {
	const op = "IsExistsById"

	row := s.db.QueryRow("SELECT id FROM users WHERE id = ?", id)
	if err := row.Err(); err != nil {
		util.LogError(f, op, err)

		return true
	}

	userItem := entity.User{}

	err := row.Scan(&userItem.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}

		util.LogError(f, op, err)
	}

	return true
}

func (s *Storage) Delete(id int64) {
	const op = "Delete"

	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		util.LogError(f, op, err)

		return
	}
}
