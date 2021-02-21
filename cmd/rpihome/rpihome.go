package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"github.com/vasilpatelnya/rpi-home/model"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "c", "config/development.json", "config path")
}

func main() {
	flag.Parse()
	log.Println("The application was launched to the path to the configuration file:", configPath)

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
	go apiServer(appContainer.DB.Mongo)
	appContainer.Run()
}

func apiServer(mongo *config.MongoConnection) {
	http.HandleFunc("/api/v1/motioneye", func(w http.ResponseWriter, r *http.Request) {
		type DetectRequest struct {
			Device string `json:"device"`
			Type   int    `json:"type"`
		}
		var request DetectRequest
		if err := JsonDecode(r.Body, &request); err != nil {
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

		err := mongo.C("events").Insert(event) // todo to cfg
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Создана запись в БД")
	})
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Api server error: %s", err.Error())
	}
}

func JsonDecode(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)

	return err
}
