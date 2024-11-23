package weather

import (
	"encoding/json"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"math"
	"testing"
)

type mockAPIClient struct {
	ResponseData []byte
	ResponseErr  error
}

func (m *mockAPIClient) DoRequest(params string) ([]byte, error) {
	return m.ResponseData, m.ResponseErr
}

func TestGetByCityName(t *testing.T) {
	mockResponse := entity.OpenWeather{
		Name: "Moscow",
		Coord: entity.Coord{
			Lon: 37.615600,
			Lat: 55.752200,
		},
		Main: entity.Main{
			Temp: 0.96,
		},
		Weather: []entity.Weather{{Icon: "10d"}},
	}

	mockData, _ := json.Marshal(mockResponse)

	s := NewService(&mockAPIClient{
		ResponseData: mockData,
		ResponseErr:  nil,
	})

	item := s.GetByCityName("Moscow")
	if item.City != mockResponse.Name {
		t.Errorf("GetByCityName: expected name %s, got %s", mockResponse.Name, item.City)
	}

	if item.Icon != mockResponse.Weather[0].Icon {
		t.Errorf("GetByCityName: expected icon %s, got %s", mockResponse.Weather[0].Icon, item.Icon)
	}
}

func TestGetByCoordinates(t *testing.T) {
	mockResponse := entity.OpenWeather{
		Name: "Moscow",
		Coord: entity.Coord{
			Lon: 37.615600,
			Lat: 55.752200,
		},
		Main: entity.Main{
			Temp: 0.96,
		},
		Weather: []entity.Weather{{Icon: "10d"}},
	}

	mockData, _ := json.Marshal(mockResponse)

	s := NewService(&mockAPIClient{
		ResponseData: mockData,
		ResponseErr:  nil,
	})

	item := s.GetByCoordinates("55.752200", "37.615600")
	if item.City != mockResponse.Name {
		t.Errorf("GetByCoordinates: expected name %s, got %s", mockResponse.Name, item.City)
	}

	if item.Temperature != int64(math.Round(mockResponse.Main.Temp)) {
		t.Errorf(
			"GetByCoordinates: expected temperature %d, got %d",
			int64(math.Round(mockResponse.Main.Temp)),
			item.Temperature,
		)
	}

	if item.Longitude != mockResponse.Coord.Lon {
		t.Errorf("GetByCoordinates: expected longitude %f, got %f", mockResponse.Coord.Lon, item.Longitude)
	}

	if item.Latitude != mockResponse.Coord.Lat {
		t.Errorf("GetByCoordinates: expected latitude %f, got %f", mockResponse.Coord.Lat, item.Latitude)
	}
}
