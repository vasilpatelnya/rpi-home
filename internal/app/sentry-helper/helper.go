package sentry_helper

import (
	"github.com/getsentry/sentry-go"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"log"
	"os"
)

func Start() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://57b796bab50e4584abc4dd1bc02e0afd@o431527.ingest.sentry.io/5382874",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	} else {
		if os.Getenv("APP_MODE") == config.AppProd {
			log.Println("sentry работает")
		}
	}
}

func Handle(err error, msg string) {
	if os.Getenv("APP_MODE") == config.AppProd {
		sentry.CaptureException(err)
	}
	log.Println(msg)
}
