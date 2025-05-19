package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/denyshuzovskyi/nimbus-notify/internal/config"
	"github.com/denyshuzovskyi/nimbus-notify/internal/dto"
	"github.com/denyshuzovskyi/nimbus-notify/internal/lib/sqlutil"
	"github.com/denyshuzovskyi/nimbus-notify/internal/model"
	"log/slog"
	"time"
)

type NotificationService struct {
	db                     *sql.DB
	weatherProvider        WeatherProvider
	locationRepository     LocationRepository
	weatherRepository      WeatherRepository
	subscriberRepository   SubscriberRepository
	subscriptionRepository SubscriptionRepository
	tokenRepository        TokenRepository
	emailSender            EmailSender
	log                    *slog.Logger
}

func NewNotificationService(
	db *sql.DB,
	weatherProvider WeatherProvider,
	locationRepository LocationRepository,
	weatherRepository WeatherRepository,
	subscriberRepository SubscriberRepository,
	subscriptionRepository SubscriptionRepository,
	tokenRepository TokenRepository,
	emailSender EmailSender,
	log *slog.Logger) *NotificationService {
	return &NotificationService{
		db:                     db,
		weatherProvider:        weatherProvider,
		locationRepository:     locationRepository,
		weatherRepository:      weatherRepository,
		subscriberRepository:   subscriberRepository,
		subscriptionRepository: subscriptionRepository,
		tokenRepository:        tokenRepository,
		emailSender:            emailSender,
		log:                    log,
	}
}

func (s *NotificationService) SendDailyNotifications(emailData config.EmailData) {
	s.log.Info("triggered SendDailyNotifications")
	ctx := context.Background()

	err := sqlutil.WithTx(ctx, s.db, &sql.TxOptions{ReadOnly: true}, func(tx *sql.Tx) error {
		subscriptions, errIn := s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, tx, model.Frequency_Daily)
		if errIn != nil {
			return errIn
		}

		if errIn = s.sendNotificationsToAllSubscribers(ctx, tx, subscriptions, emailData); errIn != nil {
			return errIn
		}

		return nil
	})
	if err != nil {
		s.log.Error("rolled back transaction because of ", "error", err)
		return
	}
	s.log.Info("transaction commited successfully")
}

func (s *NotificationService) SendHourlyNotifications(emailData config.EmailData) {
	s.log.Info("triggered SendHourlyNotifications")
	ctx := context.Background()

	err := sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		subscriptions, errIn := s.subscriptionRepository.FindAllByFrequencyAndConfirmedStatus(ctx, tx, model.Frequency_Hourly)
		if errIn != nil {
			return errIn
		}

		if errIn = s.sendNotificationsToAllSubscribers(ctx, tx, subscriptions, emailData); errIn != nil {
			return errIn
		}

		return nil
	})
	if err != nil {
		s.log.Error("rolled back transaction because of ", "error", err)
		return
	}
	s.log.Info("transaction commited successfully")
}

func (s *NotificationService) sendNotificationsToAllSubscribers(ctx context.Context, tx *sql.Tx, subscriptions []*model.Subscription, emailData config.EmailData) error {
	for i := 0; i < len(subscriptions); i++ {
		subscriber, err := s.subscriberRepository.FindById(ctx, tx, subscriptions[i].SubscriberId)
		if err != nil {
			return err
		}
		location, err := s.locationRepository.FindById(ctx, tx, subscriptions[i].LocationId)
		if err != nil {
			return err
		}
		token, err := s.tokenRepository.FindBySubscriptionIdAndType(ctx, tx, subscriptions[i].Id, model.TokenType_Unsubscribe)
		if err != nil {
			return err
		}
		lastWeather, err := s.weatherRepository.FindLastUpdatedByLocation(ctx, tx, location.Name)
		if err != nil {
			return err
		}

		if lastWeather == nil || lastWeather.LastUpdated.Add(15*time.Minute).Before(time.Now()) {
			weather, err := s.weatherProvider.GetCurrentWeather(location.Name)
			if err != nil {
				return err
			}

			weather.Weather.LocationId = location.Id
			weather.Weather.FetchedAt = time.Now().UTC()

			err = s.weatherRepository.Save(ctx, tx, &weather.Weather)
			if err != nil {
				return err
			}

			lastWeather = &weather.Weather
		}

		email := dto.SimpleEmail{
			From:    emailData.From,
			To:      subscriber.Email,
			Subject: emailData.Subject,
			Text: fmt.Sprintf(
				emailData.Text,
				location.Name,
				lastWeather.Temperature,
				lastWeather.Humidity,
				lastWeather.Description,
				token.Token,
			),
		}

		err = s.emailSender.Send(ctx, email)
		if err != nil {
			return err
		}
		s.log.Info("weather email is send")
	}

	return nil
}
