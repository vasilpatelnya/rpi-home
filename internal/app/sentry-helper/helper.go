package sentry_helper

import (
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func Start(logger *logrus.Logger, url string) {
	if url == "" {
		logger.Warning("Не указан URL для Sentry. Сообщения логгера не идут в Sentry.")

		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: url,
	})
	if err != nil {
		logger.Errorf("sentry.Init: %s", err)
	} else {
		logger.Info("Служба логгирования в Sentry успешно инициализирована.")
	}
}

func Handle(logger *logrus.Logger, err error, msg string) {
	sentry.CaptureException(err)
	logger.Debug(msg)
}
