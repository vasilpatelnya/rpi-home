package main

import (
	"flag"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"github.com/vasilpatelnya/rpi-home/internal/app/store"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"time"

	rpidetectormongo "github.com/vasilpatelnya/rpi-home/internal/app/rpi-detector-mongo"
)

var defaultDevice = "'неизвестное имя'"

func main() {
	var event = rpidetectormongo.New()
	c, err := config.New("configs/dev.env")
	if err != nil {
		log.Fatal("Ошибка загрузки env файла.")
	}
	s, err := store.New(c)
	if err != nil {
		log.Fatal("Ошибка создания подключения к БД.")
	}
	flag.StringVar(&event.Device, "device", defaultDevice, "Имя камеры")
	flag.IntVar(&event.Type, "type", rpidetectormongo.TypeUndefined, "Тип события")
	flag.Parse()
	log.Printf("Камера %s обнаружила движение. Тип события: %d", event.Device, event.Type)
	if event.Name != defaultDevice && event.Type != rpidetectormongo.TypeUndefined {
		event.ID = bson.NewObjectId()
		event.Name = "Обнаружено движение!"
		if event.Type == rpidetectormongo.TypeMovieReady {
			event.Name = "Новое видео готово!"
		}
		event.Created = time.Now().UnixNano()
		event.Updated = time.Now().UnixNano()

		err = s.Collection.Insert(event)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Создана запись в БД")

		return
	}
	log.Println("Не указан тип события или название устройства.")
	os.Exit(1)
}
