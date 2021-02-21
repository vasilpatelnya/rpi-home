package servicecontainer

import (
	"fmt"
	"github.com/vasilpatelnya/rpi-home/model"
	"github.com/vasilpatelnya/rpi-home/tool/jsontool"
	"log"
	"net/http"
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

		log.Printf("Request successfully decoded: device '%s', type '%d'", request.Device, request.Type)

		event := model.New()
		event.Device = request.Device
		event.Type = request.Type
		event.Name = "Обнаружено движение!"
		if event.Type == model.TypeMovieReady {
			event.Name = "Новое видео готово!"
		}

		event.Created = time.Now().UnixNano()
		event.Updated = time.Now().UnixNano()

		err := sc.DB.Mongo.C("events").Insert(event) // todo to cfg
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
				sc.Logger.Errorf("apiserver error: %s", apiErr.Error())
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
		options := sc.AppConfig.Notifier.Options
		sc.Notifier = telegram.New(options.Token, options.ChatID)
	}

	return nil
}

// Run ...
func (sc *ServiceContainer) Run() {
	mainTicker := time.NewTicker(sc.AppConfig.Periods.MainTickerTime * time.Millisecond)
	defer mainTicker.Stop()

	sentryhelper.Start(sc.Logger, sc.AppConfig.SentrySettings.SentryUrl)

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
