package servicecontainer

import (
	"fmt"
	"log"
	"time"

	"github.com/vasilpatelnya/rpi-home/container/apiserver"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/sqlite3"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
	"github.com/vasilpatelnya/rpi-home/usecase"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
)

// ServiceContainer ...
type ServiceContainer struct {
	AppConfig *config.Config
	DB        *dataservice.ConnectionContainer
	Logger    *logrus.Logger
	Notifier  notification.Notifier
	Repo      dataservice.EventData
}

func text(code int) string {
	return translate.T().Text(code)
}

func scErrorMsg(code int, err error) string {
	return fmt.Sprintf("%s: %s", text(code), err.Error())
}

func scErrWrap(code int, err error) error {
	msg := text(code)
	return errors.New(fmt.Sprintf("%s: %s", msg, err.Error()))
}

// InitApp initializes container config in the specified path.
func (sc *ServiceContainer) InitApp() error {
	envMode, err := config.ParseEnvMode()
	if err != nil {
		log.Fatal(scErrorMsg(translate.ErrorParsingEnv, err))
	}

	rootPath, err := fs.RootPath()
	if err != nil {
		sc.Logger.Fatal(scErrorMsg(translate.ErrorRootPath, err))
	}

	sc.AppConfig, err = config.New(fmt.Sprintf("%s/config/%s.json", rootPath, envMode))
	if err != nil {
		return scErrWrap(translate.ErrorConfigLoad, err)
	}

	err = sc.InitLogger()
	if err != nil {
		return scErrWrap(translate.ErrorParsingEnv, err)
	}

	sc.DB, err = dataservice.AssertCreateConnectionContainer(sc.AppConfig.Database)
	if err != nil {
		return scErrWrap(translate.ErrorCreateConnectionContainer, err)
	}

	if sc.AppConfig.Notifier.IsUsing {
		err = sc.InitNotifier()
		if err != nil {
			return scErrWrap(translate.ErrorNotifierInit, err)
		}
	}

	err = sc.InitRepo()
	if err != nil {
		return scErrWrap(translate.ErrorRepoInit, err)
	}

	go sc.InitApiServer()

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

// InitApiServer ...
func (sc *ServiceContainer) InitApiServer() {
	s := apiserver.New(&apiserver.ApiOpts{
		Settings: sc.AppConfig.ApiServer,
		Repo:     sc.Repo,
		Logger:   sc.Logger,
		Notifier: sc.Notifier,
	})
	s.Run()
}

// InitNotifier ...
func (sc *ServiceContainer) InitNotifier() error {
	switch sc.AppConfig.Notifier.Type {
	case notification.TypeTelegram:
		var err error
		options := sc.AppConfig.Notifier.Options
		sc.Notifier, err = telegram.New(options.Token, options.ChatID)
		if err != nil {
			return err
		}
	default:
		return errors.New("Unknown notifier type: " + sc.AppConfig.Notifier.Type)
	}

	return nil
}

// Run ...
func (sc *ServiceContainer) Run() {
	mainTicker := time.NewTicker(sc.AppConfig.Periods.MainTickerTime * time.Millisecond)
	defer mainTicker.Stop()

	sentryhelper.Start(sc.Logger, sc.AppConfig.SentrySettings.SentryUrl)

	rootPath, err := fs.RootPath()
	if err != nil {
		sc.Logger.Fatalf("root path not founded: %s", err.Error())
	}

	for {
		select {
		case <-mainTicker.C:
			opts := usecase.EventHandleOpts{
				TargetDir:   sc.AppConfig.Motion.MoviesDirCam1,
				BackupDir:   rootPath + "/backup",
				Ext:         sc.AppConfig.Motion.FileExtension,
				Repo:        sc.Repo,
				Notifier:    sc.Notifier,
				Logger:      sc.Logger,
				UseNotifier: sc.AppConfig.Notifier.IsUsing,
			}
			usecase.EventHandle(opts)
		}
	}
}

// InitRepo ...
func (sc *ServiceContainer) InitRepo() error {
	switch {
	case sc.DB.Mongo != nil:
		sc.Repo = GetRepo(sc.DB.Mongo, sc.Logger)

		return nil
	case sc.DB.SQLite3 != nil:
		sc.Repo = GetRepo(sc.DB.SQLite3, sc.Logger)

		return nil
	default:
		return errors.New("not found db connection")
	}
}

func GetRepo(connection interface{}, logger *logrus.Logger) dataservice.EventData {
	mongoConnection, isMongoConnection := connection.(*mongodb.MongoConnection)
	sqlite3Connection, isSQLite3Connection := connection.(*sqlite3.SQLite3Connection)
	switch {
	case isMongoConnection:
		return &mongodb.EventDataMongo{
			EventsCollection: mongoConnection.C("events"), // todo to cfg
			Logger:           logger,
		}
	case isSQLite3Connection:
		_, db := sqlite3Connection.C()
		return &sqlite3.EventDataSQLite3{
			DB:     db,
			Logger: logger,
		}
	default:
		logger.Error("Unknown connection")

		return nil
	}
}
