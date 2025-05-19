package service

import (
	"context"
	"github.com/denyshuzovskyi/nimbus-notify/internal/dto"
	"github.com/denyshuzovskyi/nimbus-notify/internal/lib/sqlutil"
	"github.com/denyshuzovskyi/nimbus-notify/internal/model"
)

type WeatherProvider interface {
	GetCurrentWeather(string) (*model.WeatherWithLocation, error)
}

type LocationRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Location) (int32, error)
	FindByName(context.Context, sqlutil.SQLExecutor, string) (*model.Location, error)
	FindById(context.Context, sqlutil.SQLExecutor, int32) (*model.Location, error)
}

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}
