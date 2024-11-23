package storage

import (
	"database/sql"
	"github.com/albakov/go-weather-viewer/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func MustNew(config *config.Config) *sql.DB {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)

	return db
}
