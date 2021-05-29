package testhelpers

import (
	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"path"
	"runtime"
	"strings"
)

// GetTestContainer ...
func GetTestContainer() (*servicecontainer.ServiceContainer, error) {
	appConfig := config.Config{}
	c := servicecontainer.ServiceContainer{AppConfig: &appConfig}

	err := c.InitApp()
	if err != nil {
		return nil, errors.Wrap(err, "create test container")
	}

	return &c, nil
}

func GetTestDataDir() string {
	_, thisFilename, _, _ := runtime.Caller(0)

	return strings.Replace(path.Dir(thisFilename), "tool/testhelpers", "testdata", -1)
}
