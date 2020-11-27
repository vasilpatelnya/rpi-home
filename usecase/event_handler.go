package usecase

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/internal/app/tgpost"
	"github.com/vasilpatelnya/rpi-home/model"
)

const (
	MaxSize          int64 = 50 * 1024 * 1024 // 50mb
	StatusUpdated          = 11
	StatusNotUpdated       = -11
)

func EventHandle(logger *logrus.Logger, repo dataservice.EventData, moviesPath string) {
	handleNew(logger, repo, moviesPath)
	handleFail(logger, repo, moviesPath)
}

func handleNew(logger *logrus.Logger, repo dataservice.EventData, moviesPath string) {
	handleEvents(logger, repo, model.StatusNew, moviesPath, "./backup") // todo to cfg
}

func handleFail(logger *logrus.Logger, repo dataservice.EventData, moviesPath string) {
	handleEvents(logger, repo, model.StatusFail, moviesPath, "./backup") // todo to cfg
}

func handleEvents(logger *logrus.Logger, repo dataservice.EventData, status int, moviesPath, backupPath string) {
	events, err := repo.GetAllByStatus(status)
	if err != nil {
		sentryhelper.Handle(logger, err, "Ошибка получения записей событий из БД")
	}

	if len(events) > 0 {
		for _, e := range events {
			switch e.Type {
			case model.TypeMotion:
				status, err := handleMotionAlarm(repo, &e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(logger, err, msg)
					if status == tgpost.StatusNotSent {
						e.Status = model.StatusFail
						err = repo.Save(&e)
						if err != nil {
							msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
							sentryhelper.Handle(logger, err, msg)
						}
						continue
					}
				}
				e.Status = model.StatusReady
				err = repo.Save(&e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(logger, err, msg)
				}
			case model.TypeMovieReady:
				logger.Info("Видео готово!")
				e.Status, err = handleMotionReady(&e, logger, moviesPath, backupPath)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(logger, err, msg)
				}

				err = repo.SaveUpdated(&e, e.Status)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(logger, err, msg)
				}
			}
		}
	}
}

func handleMotionAlarm(repo dataservice.EventData, e *model.Event) (int, error) {
	err := tgpost.SendText(e.GetMotionMessage())
	if err != nil {
		return tgpost.StatusNotSent, errors.New("ошибка отправки текста о срабатывании")
	}
	err = repo.Save(e)
	if err != nil {
		return StatusNotUpdated, errors.New(fmt.Sprintf("ошибка обновления записи в БД, id записи: %s", e.ID.Hex()))
	}

	return StatusUpdated, nil
}

func handleMotionReady(e *model.Event, logger *logrus.Logger, dirname string, backupPath string) (int, error) {
	l, err := tgpost.GetTodayFileList(dirname)
	if err != nil {
		logger.Error("Ошибка получения списка файлов в директории:", err.Error())

		return tgpost.StatusNotSent, err
	}
	for _, f := range l {
		todayDir := tgpost.GetTodayDir()
		fp := fmt.Sprintf("%s/%s/%s", dirname, todayDir, f.Name())
		ext := filepath.Ext(f.Name())
		if ext == ".mp4" && f.Size() > 0 { // todo to cfg
			if f.Size() < MaxSize {
				if os.Getenv("APP_MODE") != "test" {
					msg := e.GetVideoReadyMessage()
					err := tgpost.SendFile(fp, msg)
					if err != nil {
						logger.Error("Ошибка при попытке отправить видео", f.Name(), err)

						return tgpost.StatusNotSent, err
					}
				} else {
					logger.Info("Вы находитесь в тестовом режиме. Отправка файлов игнорируется.")
				}
				logger.Infof("файл %s был отправлен в телеграм", fp)
				box, err := ioutil.ReadFile(fp)
				if err != nil {
					logger.Error("Ошибка при попытке прочитать файл:", f.Name(), err)

					return tgpost.StatusNotSent, err
				}
				err = ioutil.WriteFile(backupPath+"/"+f.Name(), box, 0777)
				if err != nil {
					logger.Error("Ошибка при попытке скопировать файл:", f.Name(), err)

					return tgpost.StatusNotSent, err
				}
				err = os.Remove(fp)
				if err != nil {
					logger.Error("Ошибка при попытке удалить файл:", f.Name(), err)

					return tgpost.StatusNotSent, err
				}
			} else {
				// TODO чтобы постоянно не отсылать сообщение надо где-то зафиксировать отправку сообщения
				if os.Getenv("APP_MODE") != "test" && os.Getenv("APP_MODE") != "prod" {
					err := tgpost.SendText("Файл слишком велик чтобы его пересылать в Telegram. Вы можете его посмотреть через веб-интерфейс. Имя файла: " + f.Name())
					sentryhelper.Handle(logger, err, "Не удалось отправить текстовое сообщение о превышении размера видеофайла.")
				}
			}
		}
	}

	return tgpost.StatusSent, nil
}
