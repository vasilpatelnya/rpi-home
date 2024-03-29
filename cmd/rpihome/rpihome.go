package main

import (
	"log"

	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
)

func main() {
	appContainer, err := buildContainer()
	if err != nil {
		log.Fatal(translate.T().Text(translate.ErrorCreateContainer), err.Error())
	}
	appContainer.Run()
}

func buildContainer() (*servicecontainer.ServiceContainer, error) {
	c := &servicecontainer.ServiceContainer{AppConfig: &config.Config{}}

	return c, c.InitApp()
}
