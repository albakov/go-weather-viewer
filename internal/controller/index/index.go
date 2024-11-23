package index

import (
	"database/sql"
	"fmt"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/services/weather"
	"github.com/albakov/go-weather-viewer/internal/storage/location"
	"net/http"
)

type Controller struct {
	commonController *controller.Controller
	locationStorage  *location.Storage
	weatherService   *weather.Service
}

func New(commonController *controller.Controller, config *config.Config, db *sql.DB) *Controller {
	return &Controller{
		commonController: commonController,
		locationStorage:  location.NewStorage(db),
		weatherService:   weather.NewService(weather.NewAPI(config)),
	}
}

func (cc *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		cc.commonController.ShowNotFound(w, r)

		return
	}

	if r.Method != http.MethodGet {
		cc.commonController.ShowMethodNotAllowedError(w, r)

		return
	}

	header := cc.commonController.GetHeaderData(r)

	cc.commonController.ShowResponse(
		w,
		"index.html",
		entity.PageData{
			PageTitle: "Home",
			Header:    header,
			Data:      cc.pageData(header.User.Id),
		},
	)
}

func (cc *Controller) pageData(userId int64) entity.Index {
	locations := []entity.LocationWithTemperature{}

	if userId != 0 {
		items := cc.locationStorage.GetByUserId(userId)

		for _, item := range items {
			w := cc.weatherService.GetByCoordinates(fmt.Sprintf("%f", item.Latitude), fmt.Sprintf("%f", item.Longitude))

			locations = append(locations, entity.LocationWithTemperature{
				LocationId:  item.Id,
				Name:        item.Name,
				Temperature: w.Temperature,
				Icon:        w.Icon,
				Latitude:    w.Latitude,
				Longitude:   w.Longitude,
			})
		}
	}

	return entity.Index{
		Locations:       locations,
		HasLocations:    len(locations) > 0,
		IsAuthenticated: userId != 0,
	}
}
