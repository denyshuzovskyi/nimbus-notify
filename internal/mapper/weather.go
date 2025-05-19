package mapper

import (
	"github.com/denyshuzovskyi/nimbus-notify/internal/dto"
	"github.com/denyshuzovskyi/nimbus-notify/internal/model"
)

func WeatherToWeatherDTO(weather model.Weather) dto.WeatherDTO {
	return dto.WeatherDTO{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}
}
