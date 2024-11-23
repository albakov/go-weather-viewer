package weather

import (
	"encoding/json"
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/util"
	"math"
	"net/url"
)

const f = "weather.Service"

type Service struct {
	apiClient APIClient
}

func NewService(apiClient APIClient) *Service {
	return &Service{apiClient: apiClient}
}

func (s *Service) GetByCityName(city string) entity.WeatherItem {
	const op = "GetByCityName"

	resBody, err := s.apiClient.DoRequest(fmt.Sprintf("&q=%s", url.QueryEscape(city)))
	if err != nil {
		util.LogError(f, op, err)

		return entity.WeatherItem{}
	}

	if len(resBody) == 0 {
		return entity.WeatherItem{}
	}

	openWeather := entity.OpenWeather{}
	err = json.Unmarshal(resBody, &openWeather)
	if err != nil {
		util.LogError(f, op, err)

		return entity.WeatherItem{}
	}

	return entity.WeatherItem{
		City:        openWeather.Name,
		Temperature: int64(math.Round(openWeather.Main.Temp)),
		Latitude:    openWeather.Coord.Lat,
		Longitude:   openWeather.Coord.Lon,
		Icon:        s.getIcon(openWeather),
	}
}

func (s *Service) GetByCoordinates(latitude, longitude string) entity.WeatherItem {
	const op = "GetByCoordinates"

	resBody, err := s.apiClient.DoRequest(fmt.Sprintf(
		"&lat=%s&lon=%s",
		url.QueryEscape(latitude),
		url.QueryEscape(longitude),
	))
	if err != nil {
		util.LogError(f, op, err)

		return entity.WeatherItem{}
	}

	if len(resBody) == 0 {
		return entity.WeatherItem{}
	}

	openWeather := entity.OpenWeather{}
	err = json.Unmarshal(resBody, &openWeather)
	if err != nil {
		util.LogError(f, op, err)

		return entity.WeatherItem{}
	}

	return entity.WeatherItem{
		City:        openWeather.Name,
		Temperature: int64(math.Round(openWeather.Main.Temp)),
		Latitude:    openWeather.Coord.Lat,
		Longitude:   openWeather.Coord.Lon,
		Icon:        s.getIcon(openWeather),
	}
}

func (s *Service) getIcon(openWeather entity.OpenWeather) string {
	icon := ""

	if len(openWeather.Weather) > 0 {
		icon = openWeather.Weather[0].Icon
	}

	return icon
}
