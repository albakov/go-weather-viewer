package weather

import (
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/util"
	"io"
	"net/http"
	"time"
)

type APIClient interface {
	DoRequest(params string) ([]byte, error)
}

type API struct {
	config *config.Config
	client *http.Client
}

func NewAPI(config *config.Config) *API {
	return &API{
		config: config,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (h *API) DoRequest(params string) ([]byte, error) {
	const f = "weather.APIClient"
	const op = "DoRequest"

	res, err := h.client.Get(fmt.Sprintf(
		"%s?appid=%s&lang=ru&units=metric%s",
		h.config.WeatherApi.Url,
		h.config.WeatherApi.Key,
		params,
	))
	if err != nil {
		util.LogError(f, op, err)

		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		util.LogError(f, op, err)

		return nil, err
	}
	defer res.Body.Close()

	return resBody, nil
}
