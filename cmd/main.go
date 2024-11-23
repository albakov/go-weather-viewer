package main

import (
	"github.com/albakov/go-weather-viewer/internal/app"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/storage"
	"log"
	"os"
	"path/filepath"
)

func main() {
	c := config.MustNew(getConfigPath())
	db := storage.MustNew(c)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	app.New(c, db).MustStart()
}

func getConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, "config", "app.toml")
}
