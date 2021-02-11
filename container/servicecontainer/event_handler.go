package usecase

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/model"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
)

const (
	MaxSize          int64 = 50 * 1024 * 1024 // 50mb
	StatusUpdated          = 11
	StatusNotUpdated       = -11
)

// todo переделать на аргумент с несколькими путями к записям: например для нескольких камер
func EventHandle(sc *servicecontainer.ServiceContainer, repo dataservice.EventData, moviesPath string) {
	handleNew(sc, repo, moviesPath)
	handleFail(sc, repo, moviesPath)
}

func handleNew(sc *servicecontainer.ServiceContainer, repo dataservice.EventData, moviesPath string) {
	handleEvents(sc, repo, model.StatusNew, moviesPath, "./backup") // todo to cfg
}

func handleFail(sc *servicecontainer.ServiceContainer, repo dataservice.EventData, moviesPath string) {
	handleEvents(sc, repo, model.StatusFail, moviesPath, "./backup") // todo to cfg
}

func handleEvents(sc *servicecontainer.ServiceContainer, repo dataservice.EventData, status int, moviesPath, backupPath string) {
	events, err := repo.GetAllByStatus(status)
	if err != nil {
		sentryhelper.Handle(sc.Logger, err, "Ошибка получения записей событий из БД")
	}

	if len(events) > 0 {
		for _, e := range events {
			switch e.Type {
			case model.TypeMotion:
				status, err := handleMotionAlarm(sc.Notifier, repo, &e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
					if status == model.StatusNotSent {
						e.Status = model.StatusFail
						err = repo.Save(&e)
						if err != nil {
							msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
							sentryhelper.Handle(sc.Logger, err, msg)
						}
						continue
					}
				}
				e.Status = model.StatusReady
				err = repo.Save(&e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
				}
			case model.TypeMovieReady:
				sc.Logger.Info("Видео готово!")
				e.Status, err = handleMotionReady(sc, &e, moviesPath, backupPath)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
				}

				err = repo.SaveUpdated(&e, e.Status)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
				}
			}
		}
	}
}

func handleMotionAlarm(notifier notification.Notifier, repo dataservice.EventData, e *model.Event) (int, error) {
	err := notifier.SendText(e.GetMotionMessage())
	if err != nil {
		return model.StatusNotSent, errors.New("ошибка отправки текста о срабатывании")
	}
	err = repo.Save(e)
	if err != nil {
		return StatusNotUpdated, errors.New(fmt.Sprintf("ошибка обновления записи в БД, id записи: %s", e.ID.Hex()))
	}

	return StatusUpdated, nil
}

func handleMotionReady(sc *servicecontainer.ServiceContainer, e *model.Event, dirname string, backupPath string) (int, error) {
	l, err := fs.GetTodayFileList(dirname, model.LayoutISO)
	if err != nil {
		sc.Logger.Error("Ошибка получения списка файлов в директории:", err.Error())

		return model.StatusNotSent, err
	}
	for _, f := range l {
		todayDir := fs.GetTodayDir(model.LayoutISO)
		fp := fmt.Sprintf("%s/%s/%s", dirname, todayDir, f.Name())
		ext := filepath.Ext(f.Name())
		if ext == ".mp4" && f.Size() > 0 { // todo to cfg
			if f.Size() < MaxSize {
				if os.Getenv("APP_MODE") != "test" {
					msg := e.GetVideoReadyMessage()
					err := sc.Notifier.SendFile(fp, msg)
					if err != nil {
						sc.Logger.Error("Ошибка при попытке отправить видео", f.Name(), err)

						return model.StatusNotSent, err
					}
				} else {
					sc.Logger.Info("Вы находитесь в тестовом режиме. Отправка файлов игнорируется.")
				}
				sc.Logger.Infof("файл %s был отправлен в телеграм", fp)
				box, err := ioutil.ReadFile(fp)
				if err != nil {
					sc.Logger.Error("Ошибка при попытке прочитать файл:", f.Name(), err)

					return model.StatusNotSent, err
				}
				err = ioutil.WriteFile(backupPath+"/"+f.Name(), box, 0777)
				if err != nil {
					sc.Logger.Error("Ошибка при попытке скопировать файл:", f.Name(), err)

					return model.StatusNotSent, err
				}
				err = os.Remove(fp)
				if err != nil {
					sc.Logger.Error("Ошибка при попытке удалить файл:", f.Name(), err)

					return model.StatusNotSent, err
				}
			} else {
				// TODO чтобы постоянно не отсылать сообщение надо где-то зафиксировать отправку сообщения
				if os.Getenv("APP_MODE") != "test" && os.Getenv("APP_MODE") != "prod" {
					err := sc.Notifier.SendText("Файл слишком велик чтобы его пересылать в Telegram. Вы можете его посмотреть через веб-интерфейс. Имя файла: " + f.Name())
					sentryhelper.Handle(sc.Logger, err, "Не удалось отправить текстовое сообщение о превышении размера видеофайла.")
				}
			}
		}
	}

	return model.StatusSent, nil
}