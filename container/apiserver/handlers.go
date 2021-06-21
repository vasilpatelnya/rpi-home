package apiserver

import (
	"fmt"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/model"
	"github.com/vasilpatelnya/rpi-home/tool/jsontool"
	"log"
	"net/http"
	"time"
)

func MotionEyeHandler(repo dataservice.EventData) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type DetectRequest struct {
			Device string `json:"device"`
			Type   int    `json:"type"`
		}
		var request DetectRequest
		if err := jsontool.JsonDecode(r.Body, &request); err != nil {
			log.Printf("json decode error: %s\n", err.Error())
		}

		if request.Device == "" || request.Type == model.TypeUndefined {
			_, err := fmt.Fprintln(w, "Не указан тип события или название устройства.")
			if err != nil {
				log.Printf("Fprintln() error: %s", err.Error())
			}

			return
		}

		log.Printf("Request successfully decoded: device '%s', type '%d'\n", request.Device, request.Type)

		event := model.New()
		event.Device = request.Device
		event.Type = request.Type
		event.Name = "Обнаружено движение!"
		if event.Type == model.TypeMovieReady {
			event.Name = "Новое видео готово!"
		}

		event.Created = time.Now().UnixNano()
		event.Updated = time.Now().UnixNano()

		err := repo.Save(event)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Создана запись в БД")
	}
}
