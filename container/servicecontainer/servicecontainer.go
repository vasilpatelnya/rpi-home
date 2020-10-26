package servicecontainer

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/config"
)

// ServiceContainer ...
type ServiceContainer struct {
	AppConfig *config.Config
	DB        *config.ConnectionContainer
	Logger    *logrus.Logger
}

// InitApp initializes container config in the specified path.
func (sc *ServiceContainer) InitApp(filename string) error {
	c, err := config.New(filename)
	if err != nil {
		return errors.Wrap(err, "loadConfig")
	}
	sc.AppConfig = c
	sc.DB = sc.AppConfig.AssertCreateConnectionContainer()

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
