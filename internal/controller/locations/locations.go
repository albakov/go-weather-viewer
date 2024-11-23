package locations

import (
	"database/sql"
	"github.com/albakov/go-weather-viewer/internal/config"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/entity"
	"github.com/albakov/go-weather-viewer/internal/services/weather"
	"github.com/albakov/go-weather-viewer/internal/storage/location"
	locationvalidator "github.com/albakov/go-weather-viewer/internal/validation/location"
	"net/http"
	"strconv"
)

type Controller struct {
	commonController *controller.Controller
	config           *config.Config
	locationStorage  *location.Storage
	weatherService   *weather.Service
}

func New(commonController *controller.Controller, config *config.Config, db *sql.DB) *Controller {
	return &Controller{
		commonController: commonController,
		config:           config,
		locationStorage:  location.NewStorage(db),
		weatherService:   weather.NewService(weather.NewAPI(config)),
	}
}

func (cc *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cc.searchHandler(w, r)

		return
	}

	if r.Method == http.MethodPost {
		cc.createHandler(w, r)

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w, r)
}

func (cc *Controller) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		locationId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			cc.commonController.ShowNotFound(w, r)

			return
		}

		data := cc.commonController.GetHeaderData(r)
		if data.User.Id == 0 {
			cc.commonController.ShowNotFound(w, r)

			return
		}

		err = cc.locationStorage.Delete(locationId, data.User.Id)
		if err != nil {
			cc.commonController.ShowNotFound(w, r)

			return
		}

		cc.commonController.RedirectTo(w, "/")

		return
	}

	cc.commonController.ShowMethodNotAllowedError(w, r)
}

func (cc *Controller) searchHandler(w http.ResponseWriter, r *http.Request) {
	weatherItem := entity.WeatherItem{}
	requestedCity := r.URL.Query().Get("city")
	if requestedCity != "" {
		weatherItem = cc.weatherService.GetByCityName(requestedCity)
	}

	cc.commonController.ShowResponse(
		w,
		"locations.html",
		entity.PageData{
			PageTitle: "Locations",
			Header:    cc.commonController.GetHeaderData(r),
			Data: entity.LocationPage{
				WeatherItem:     weatherItem,
				HasWeatherItems: weatherItem.City != "",
				SearchedCity:    requestedCity,
			},
		},
	)
}

func (cc *Controller) createHandler(w http.ResponseWriter, r *http.Request) {
	data := cc.commonController.GetHeaderData(r)

	if data.User.Id == 0 {
		cc.commonController.RedirectTo(w, "/login")

		return
	}

	validator := locationvalidator.NewValidator(r, cc.config)
	validator.Validate()

	if !validator.IsValid() {
		cc.commonController.ShowNotFound(w, r)

		return
	}

	if cc.locationStorage.ExistsByUserIdAndName(data.User.Id, validator.ValueByName("name")) {
		cc.commonController.RedirectTo(w, "/")

		return
	}

	err := cc.locationStorage.Create(entity.Location{
		Name:      validator.ValueByName("name"),
		UserId:    data.User.Id,
		Latitude:  validator.GetFloat("latitude"),
		Longitude: validator.GetFloat("longitude"),
	})
	if err != nil {
		cc.commonController.ShowNotFound(w, r)

		return
	}

	cc.commonController.RedirectTo(w, "/")
}
