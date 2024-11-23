package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Host                string `toml:"host"`
	Port                int64  `toml:"port"`
	DSN                 string `toml:"dsn"`
	DSNTest             string `toml:"dsn_test"`
	SessionPeriod       int64  `toml:"session_period"`
	FieldValueMinLength int64  `toml:"field_value_min_length"`
	WeatherApi
}

type WeatherApi struct {
	Url string `toml:"weather_api_url"`
	Key string `toml:"weather_api_key"`
}

func MustNew(pathToConfig string) *Config {
	c := &Config{}

	_, err := toml.DecodeFile(pathToConfig, c)
	if err != nil {
		panic(err)
	}

	return c
}
