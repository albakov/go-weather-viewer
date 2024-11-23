package location

import (
	"database/sql"
	"errors"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/util"
)

const f = "location.Storage"

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(item entity.Location) error {
	const op = "Create"

	stmt, err := s.db.Prepare("INSERT INTO locations (name, user_id, latitude, longitude) VALUES (?, ?, ?, ?)")
	if err != nil {
		util.LogError(f, op, err)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Name, item.UserId, item.Latitude, item.Longitude)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	return nil
}

func (s *Storage) Delete(id, userId int64) error {
	const op = "Delete"

	stmt, err := s.db.Prepare("DELETE FROM locations WHERE id = ? AND user_id = ?")
	if err != nil {
		util.LogError(f, op, err)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, userId)
	if err != nil {
		util.LogError(f, op, err)

		return err
	}

	return nil
}

func (s *Storage) GetByUserId(userId int64) []entity.Location {
	const op = "GetByUserId"

	rows, err := s.db.Query(
		"SELECT id, name, user_id, latitude, longitude FROM locations WHERE user_id = ? LIMIT 10",
		userId,
	)
	if err != nil {
		util.LogError(f, op, err)

		return nil
	}
	defer rows.Close()

	var locations []entity.Location

	for rows.Next() {
		var item entity.Location

		err := rows.Scan(&item.Id, &item.Name, &item.UserId, &item.Latitude, &item.Longitude)
		if err != nil {
			util.LogError(f, op, err)

			return nil
		}

		locations = append(locations, item)
	}

	return locations
}

func (s *Storage) ExistsByUserIdAndName(userId int64, name string) bool {
	const op = "ExistsByUserIdAndName"

	query := s.db.QueryRow("SELECT id FROM locations WHERE user_id = ? AND name = ?", userId, name)
	if query.Err() != nil {
		util.LogError(f, op, query.Err())

		return false
	}

	item := entity.Location{}

	err := query.Scan(&item.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}

		util.LogError(f, op, err)

		return false
	}

	return item.Id != 0
}
