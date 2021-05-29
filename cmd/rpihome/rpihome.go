package main

import (
	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"log"
)

func main() {
	appContainer, err := buildContainer()
	if err != nil {
		log.Fatal("Error on try create application container:", err.Error())
	}
	appContainer.Run()
}

func buildContainer() (*servicecontainer.ServiceContainer, error) {
	appConfig := config.Config{}
	c := &servicecontainer.ServiceContainer{AppConfig: &appConfig}

	err := c.InitApp()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return c, nil
}
