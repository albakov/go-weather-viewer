package session

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"github.com/albakov/go-weather-viewer/internal/util"
)

const f = "session.Storage"

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(item entity.Session) error {
	const op = "Create"

	stmt, err := s.db.Prepare("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		util.LogError(f, op, err)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Id, item.UserId, item.ExpiresAt)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	return nil
}

func (s *Storage) Delete(id string) error {
	const op = "Delete"

	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		util.LogError(f, op, err)

		return err

	}

	return nil
}

func (s *Storage) GetById(id string) (entity.Session, error) {
	const op = "GetById"

	row := s.db.QueryRow("SELECT id, user_id, expires_at FROM sessions WHERE id = ?", id)
	if err := row.Err(); err != nil {
		util.LogError(f, op, err)

		return entity.Session{}, err
	}

	item := entity.Session{}

	err := row.Scan(&item.Id, &item.UserId, &item.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Session{}, storage.EntitiesNotFoundError
		}

		util.LogError(f, op, err)

		return entity.Session{}, err
	}

	return item, nil
}
