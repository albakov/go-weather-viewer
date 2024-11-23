package entity

type WeatherItem struct {
	City        string
	Temperature int64
	Latitude    float64
	Longitude   float64
	Icon        string
}

type OpenWeather struct {
	Name    string    `json:"name"`
	Coord   Coord     `json:"coord"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Weather struct {
	Icon string `json:"icon"`
}
