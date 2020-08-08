package sentry_helper

import (
	"github.com/getsentry/sentry-go"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"log"
	"os"
)

func Start() {
	if os.Getenv("SENTRY_URL") == "" {
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_URL"),
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
	if os.Getenv("APP_MODE") == config.AppProd && os.Getenv("SENTRY_URL") != "" {
		sentry.CaptureException(err)
	}
	log.Println(msg)
}
