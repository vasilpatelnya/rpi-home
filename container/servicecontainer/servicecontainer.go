package servicecontainer

import (
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/config"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"github.com/vasilpatelnya/rpi-home/usecase"
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
		sc.Notifier = telegram.New()
	}

	return nil
}

// Run ...
func (sc *ServiceContainer) Run() {
	timeFormat := "2 January 2006 15:04" // todo to cfg
	mainTicker := time.NewTicker(sc.AppConfig.Periods.MainTickerTime * time.Millisecond)

	sentryhelper.Start(sc.Logger, sc.AppConfig.SentrySettings.SentryUrl)

	defer mainTicker.Stop()
	for {
		select {
		case t := <-mainTicker.C:
			sc.Logger.Infof("Итерация главного цикла началась. Время: %s", t.Format(timeFormat))

			repo := &mongodb.EventDataMongo{
				EventsCollection: sc.DB.Mongo.C("events"), // todo to cfg
				Logger:           sc.Logger,
			}
			usecase.EventHandle(sc, repo, sc.AppConfig.Motion.MoviesDirCam1)

			sc.Logger.Infof("Итерация главного цикла закончилась. Время: %s", t.Format(timeFormat))
		}
	}
}
