package main

import (
	"flag"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/model"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"time"
)

var defaultDevice = "'неизвестное имя'"

func main() {
	var event = model.New()
	c, err := config.New("config/development.env") // todo fix
	if err != nil {
		log.Fatal("Ошибка загрузки env файла.")
	}
	s := c.AssertCreateConnectionContainer()

	flag.StringVar(&event.Device, "device", defaultDevice, "Имя камеры")
	flag.IntVar(&event.Type, "type", model.TypeUndefined, "Тип события")
	flag.Parse()
	log.Printf("Камера %s обнаружила движение. Тип события: %d", event.Device, event.Type)
	if event.Name != defaultDevice && event.Type != model.TypeUndefined {
		event.ID = bson.NewObjectId()
		event.Name = "Обнаружено движение!"
		if event.Type == model.TypeMovieReady {
			event.Name = "Новое видео готово!"
		}
		event.Created = time.Now().UnixNano()
		event.Updated = time.Now().UnixNano()

		err = s.Mongo.C("events").Insert(event) // todo to cfg
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Создана запись в БД")

		return
	}
	log.Println("Не указан тип события или название устройства.")
	os.Exit(1)
}
