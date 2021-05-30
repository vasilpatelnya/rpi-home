package servicecontainer

import (
	"fmt"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/sqlite3"
	"github.com/vasilpatelnya/rpi-home/model"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"github.com/vasilpatelnya/rpi-home/tool/jsontool"
	"github.com/vasilpatelnya/rpi-home/usecase"
	"log"
	"net/http"
	"time"

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

// InitApp initializes container config in the specified path.
func (sc *ServiceContainer) InitApp() error {
	envMode, err := config.ParseEnvMode()
	if err != nil {
		log.Fatalf("Parse environment mode error: %s", err.Error())
	}
	rootPath, err := fs.RootPath()
	if err != nil {
		log.Fatalf("Root path not founded, error: %s", err.Error())
	}
	sc.AppConfig, err = config.New(fmt.Sprintf("%s/config/%s.json", rootPath, envMode))
	if err != nil {
		return errors.Wrap(err, "Ошибка при загрузке конфигурационного файла:")
	}
	sc.DB, err = dataservice.AssertCreateConnectionContainer(sc.AppConfig.Database)
	if err != nil {
		return errors.Errorf("Create connection container error: %s", err.Error())
	}
	err = sc.InitLogger()
	if err != nil {
		return errors.Wrap(err, "Ошибка при инициализации логгера")
	}
	if sc.AppConfig.Notifier.IsUsing {
		err = sc.InitNotifier()
		if err != nil {
			return errors.Wrap(err, "Ошибка при инициализации модуля отправки уведомлений")
		}
	}

	switch {
	case sc.DB.Mongo != nil:
		sc.Repo = GetRepo(sc.DB.Mongo, sc.Logger)
	case sc.DB.SQLite3 != nil:
		sc.Repo = GetRepo(sc.DB.SQLite3, sc.Logger)
	default:
		return errors.New("not found db connection")
	}

	//go sc.InitApiServer() // todo add chan for manipulations

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
	http.HandleFunc("/api/v1/motioneye", func(w http.ResponseWriter, r *http.Request) {
		type DetectRequest struct {
			Device string `json:"device"`
			Type   int    `json:"type"`
		}
		var request DetectRequest
		if err := jsontool.JsonDecode(r.Body, &request); err != nil {
			log.Printf("json decode error: %s\n", err.Error())
		}

		if request.Device == "" || request.Type == model.TypeUndefined {
			_, err := fmt.Fprintln(w, "Не указан тип события или название устройства.")
			if err != nil {
				log.Printf("Fprintln() error: %s", err.Error())
			}

			return
		}

		log.Printf("Request successfully decoded: device '%s', type '%d'\n", request.Device, request.Type)

		event := model.New()
		event.Device = request.Device
		event.Type = request.Type
		event.Name = "Обнаружено движение!"
		if event.Type == model.TypeMovieReady {
			event.Name = "Новое видео готово!"
		}

		event.Created = time.Now().UnixNano()
		event.Updated = time.Now().UnixNano()

		err := sc.Repo.Save(event)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Создана запись в БД")
	})

	apiErrChan := make(chan error, 512)

	go func(apiErrChan chan error) {
		for {
			select {
			case apiErr := <-apiErrChan:
				sentryhelper.Handle(sc.Logger, apiErr, fmt.Sprintf("apiserver error: %s", apiErr.Error()))
			}
		}
	}(apiErrChan)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		msg := fmt.Sprintf("Api server error: %s", err.Error())
		sc.Logger.Errorln(msg)

		apiErrChan <- err
	}
}

// InitNotifier ...
func (sc *ServiceContainer) InitNotifier() error {
	switch sc.AppConfig.Notifier.Type {
	case "telegram":
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
				TargetDir: sc.AppConfig.Motion.MoviesDirCam1,
				BackupDir: rootPath + "/backup",
				Ext:       sc.AppConfig.Motion.FileExtension,
				Repo:      sc.Repo,
				Notifier:  sc.Notifier,
				Logger:    sc.Logger,
			}
			usecase.EventHandle(opts)
		}
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
