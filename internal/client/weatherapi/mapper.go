package weatherapi

import (
	"github.com/denyshuzovskyi/nimbus-notify/internal/model"
	"time"
)

func CurrentWeatherToWeatherWithLocation(currentWeather CurrentWeather) model.WeatherWithLocation {
	return model.WeatherWithLocation{
		Weather: model.Weather{
			LocationId:  0,
			LastUpdated: time.Unix(currentWeather.Current.LastUpdated, 0).UTC(),
			FetchedAt:   time.Unix(0, 0),
			Temperature: currentWeather.Current.TempC,
			Humidity:    float32(currentWeather.Current.Humidity),
			Description: currentWeather.Current.Condition.Text,
		},
		Location: model.Location{
			Id:   0,
			Name: currentWeather.Location.Name,
		},
	}
}
