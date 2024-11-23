package entity

type Index struct {
	Locations                     []LocationWithTemperature
	HasLocations, IsAuthenticated bool
}

type PageData struct {
	PageTitle string
	Header    Header
	Data      any
}

type Header struct {
	User User
}

type ErrorData struct {
	Code int
}

type LocationPage struct {
	WeatherItem     WeatherItem
	SearchedCity    string
	HasWeatherItems bool
}
