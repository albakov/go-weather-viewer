package entity

type Location struct {
	Id        int64
	Name      string
	UserId    int64
	Latitude  float64
	Longitude float64
}

type LocationWithTemperature struct {
	LocationId  int64
	Name        string
	Temperature int64
	Icon        string
	Latitude    float64
	Longitude   float64
}
