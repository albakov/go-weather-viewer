# WeatherViewer на Go

Подробное описание проекта по [ссылке](https://zhukovsd.github.io/java-backend-learning-course/projects/weather-viewer/).


## Конфигурация

Все опции для конфигурирования собраны в файле `config/app_example.toml` Необходимо переименовать этот файл в `app.toml`.

## Миграции

В проекте используется база данных `MariaDB`. Необходимо создать базу данных перед миграцией.

Для запуска миграций необходимо установить модуль `goose`:
https://github.com/pressly/goose (Или вручную выполнить sql-запросы из файлов `db/migrations`)

Далее выполнить команду:

`make m_up dsn={DSN}`

Здесь `{DSN}` нужно заменить на строку вида `DB_USERNAME:DB_PASSWORD@/DB_TABLE`

## Сборка
Команда для сборки:

`make build`

Или:

`go build cmd/main.go && mv main weather_viewer`

## Запуск

`./weather_viewer`

## Тесты

Для тестов используется тестовая база данных, необходимо ее создать. И выполнить миграцию:

`make m_up dsn={DSN}`

Здесь `{DSN}` нужно заменить на строку вида `DB_USERNAME:DB_PASSWORD@/DB_TABLE_TEST`. 

В `app.toml` также указать `DSN` для тестовой базы данных.