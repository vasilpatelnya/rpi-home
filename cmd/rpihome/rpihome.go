package main

import (
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
	"log"
)

func main() {
	appContainer, err := buildContainer()
	if err != nil {
		log.Fatal(translate.T().Text(translate.ErrorCreateContainer), err.Error())
	}
	appContainer.Run()
}

func buildContainer() (*servicecontainer.ServiceContainer, error) {
	appConfig := config.Config{}
	c := &servicecontainer.ServiceContainer{AppConfig: &appConfig}

	err := c.InitApp()
	if err != nil {
		return nil, err
	}

	return c, nil
}
