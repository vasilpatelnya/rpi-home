package servicecontainer

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
)

// ServiceContainer ...
type ServiceContainer struct {
	AppConfig *config.Config
	DB        *config.ConnectionContainer
	Logger    *logrus.Logger
	Notifier  notification.Notifier
}

// InitApp initializes container config in the specified path.
func (sc *ServiceContainer) InitApp(filename string) error {
	c, err := config.New(filename)
	if err != nil {
		return errors.Wrap(err, "Ошибка при загрузке конфигурационного файла:")
	}
	sc.AppConfig = c
	sc.DB = sc.AppConfig.AssertCreateConnectionContainer()
	err = sc.InitLogger()
	if err != nil {
		return errors.Wrap(err, "Ошибка при инициализации логгера")
	}
	err = sc.InitNotifier()
	if err != nil {
		return errors.Wrap(err, "Ошибка при инициализации модуля отправки уведомлений")
	}

	return nil
}

// InitLogger ...
func (sc *ServiceContainer) InitLogger() error {
	logger := logrus.New()
	sc.Logger = logger

	level, err := logrus.ParseLevel(sc.AppConfig.Logger.LogLevel)
	if err != nil {
		return err
	}
	sc.Logger.SetLevel(level)
	sc.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})
	sc.Logger.SetReportCaller(sc.AppConfig.Logger.ShowCaller)

	return err
}

// InitNotifier ...
func (sc *ServiceContainer) InitNotifier() error {
	switch sc.AppConfig.Notifier.Type {
	case "telegram":
		options := sc.AppConfig.Notifier.Options
		sc.Notifier = telegram.New(options.Token, options.ChatID)
	}

	return nil
}

// Run ...
func (sc *ServiceContainer) Run() {
	mainTicker := time.NewTicker(sc.AppConfig.Periods.MainTickerTime * time.Millisecond)

	sentryhelper.Start(sc.Logger, sc.AppConfig.SentrySettings.SentryUrl)

	defer mainTicker.Stop()
	for {
		select {
		case <-mainTicker.C:
			repo := &mongodb.EventDataMongo{
				EventsCollection: sc.DB.Mongo.C("events"), // todo to cfg
				Logger:           sc.Logger,
			}
			sc.EventHandle(repo, sc.AppConfig.Motion.MoviesDirCam1)
		}
	}
}
