package main

import (
	"flag"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/internal/app/daemon"
	rpidetectormongo "github.com/vasilpatelnya/rpi-home/internal/app/rpi-detector-mongo"
	sentryhelper "github.com/vasilpatelnya/rpi-home/internal/app/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/internal/app/tgpost"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "configs/dev.env", "путь к конфигурационному файлу")
}

func main() {
	flag.Parse()
	d := daemon.New(configPath)
	sentryhelper.Start()
	defer log.Println("Главный цикл завершился...")
	defer d.Ticker.Stop()
	for {
		select {
		case t := <-d.Ticker.C:
			log.Println("Итерация главного цикла началась.", t)

			mainHandler(d, rpidetectormongo.StatusFail)
			mainHandler(d, rpidetectormongo.StatusNew)

			log.Println("Итерация главного цикла закончилась.", t)
		}
	}
}

func mainHandler(d *daemon.Daemon, s int) {
	events, err := rpidetectormongo.GetAllByStatus(d.Store.Collection, s)
	if err != nil {
		sentryhelper.Handle(err, "Ошибка получения записей событий из БД")
	}
	if len(events) > 0 {
		for _, e := range events {
			switch e.Type {
			case rpidetectormongo.TypeMotion:
				status, err := d.MotionHandler(&e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(err, msg)
					if status == tgpost.StatusNotSent {
						e.Status = rpidetectormongo.StatusFail
						err = e.Save(d.Store.Collection)
						if err != nil {
							msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
							sentryhelper.Handle(err, msg)
						}
						continue
					}
				}
				e.Status = rpidetectormongo.StatusReady
				err = e.Save(d.Store.Collection)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(err, msg)
				}
			case rpidetectormongo.TypeMovieReady:
				log.Println("Видео готово!")
				e.Status, err = e.HandlerMotionReady(d.Config.MoviesDirCamera1, "./backup")
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(err, msg)
				}

				err = e.SaveUpdated(d.Store.Collection, e.Status)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(err, msg)
				}
			}
		}
	}
}
