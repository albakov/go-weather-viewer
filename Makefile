dev:
	go build cmd/main.go && mv main weather_viewer
	./weather_viewer

build:
	go build cmd/main.go && mv main weather_viewer

m_up:
	goose -dir db/migrations mysql $(dsn) up

m_down:
	goose -dir db/migrations mysql $(dsn) down

m_create:
	goose -dir db/migrations mysql $(dsn) create $(name) sql
