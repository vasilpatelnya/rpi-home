package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"log"
	"net/http"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config/development.json", "config path")
}

func main() {
	flag.Parse()
	log.Println("The application was launched to the path to the configuration file:", configPath)

	go apiServer()

	run()
}

func buildContainer(filename string) (*servicecontainer.ServiceContainer, error) {
	appConfig := config.Config{}
	c := &servicecontainer.ServiceContainer{AppConfig: &appConfig}

	err := c.InitApp(filename)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return c, nil
}

func run() {
	appContainer, err := buildContainer(configPath)
	if err != nil {
		log.Println("Error on try create application container:", err.Error())
		os.Exit(1)
	}
	appContainer.Run()
}

func apiServer() {
	http.HandleFunc("/detect", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "motion detected!")
	})
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Api server error: %s", err.Error())
	}
}
