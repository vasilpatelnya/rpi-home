package servicecontainer_test

import (
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
)

const (
	rootPath = "./../../"
	configPath = "config/test.json"
	RightConfigPath = rootPath + configPath
	WrongConfigPath = "./everybody.json"
)

func getTestServiceContainer() servicecontainer.ServiceContainer {
	return servicecontainer.ServiceContainer{AppConfig: &config.Config{}}
}
