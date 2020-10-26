package testhelpers

import (
	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
)

// GetTestContainer ...
func GetTestContainer(filename string) (*servicecontainer.ServiceContainer, error) {
	appConfig := config.Config{}
	c := servicecontainer.ServiceContainer{AppConfig: &appConfig}

	err := c.InitApp(filename)
	if err != nil {
		return nil, errors.Wrap(err, "create test container")
	}

	return &c, nil
}
