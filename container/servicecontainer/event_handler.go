package servicecontainer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
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
func (sc *ServiceContainer) EventHandle() {
	sc.handleEvents(model.StatusNew, sc.AppConfig.Motion.MoviesDirCam1, "/home/pi/go/src/github.com/vasilpatelnya/rpi-home/backup")  // todo to cfg
	sc.handleEvents(model.StatusFail, sc.AppConfig.Motion.MoviesDirCam1, "/home/pi/go/src/github.com/vasilpatelnya/rpi-home/backup") // todo to cfg
}

func (sc *ServiceContainer) handleEvents(status int, moviesPath, backupPath string) {
	events, err := sc.Repo.GetAllByStatus(status)
	if err != nil {
		sentryhelper.Handle(sc.Logger, err, "Ошибка получения записей событий из БД")
	}

	if len(events) > 0 {
		for _, e := range events {
			switch e.Type {
			case model.TypeMotion:
				status, err := handleMotionAlarm(sc.Notifier, sc.Repo, &e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
					if status == model.StatusNotSent {
						e.Status = model.StatusFail
						err = sc.Repo.Save(&e)
						if err != nil {
							msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
							sentryhelper.Handle(sc.Logger, err, msg)
						}
						continue
					}
				}
				e.Status = model.StatusReady
				err = sc.Repo.Save(&e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
				}
			case model.TypeMovieReady:
				sc.Logger.Info("Видео готово!")
				e.Status, err = sc.handleMotionReady(&e, moviesPath, backupPath)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(sc.Logger, err, msg)
				}

				err = sc.Repo.SaveUpdated(&e, e.Status)
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

func (sc *ServiceContainer) handleMotionReady(e *model.Event, dirname string, backupPath string) (int, error) {
	l, err := fs.GetTodayFileList(dirname, model.LayoutISO)
	if err != nil {
		sc.Logger.Errorf("Ошибка получения списка файлов в директории %s: %s", dirname, err.Error())

		return model.StatusNotSent, err
	}
	for _, f := range l {
		todayDir := fs.GetTodayDir(model.LayoutISO)
		fp := fmt.Sprintf("%s/%s/%s", dirname, todayDir, f.Name())
		ext := filepath.Ext(f.Name())
		if ext == ".mp4" && f.Size() > 0 { // todo to cfg
			if f.Size() < MaxSize {
				msg := e.GetVideoReadyMessage()
				err = fs.CopyFile(fp, "/home/pi/go/src/github.com/vasilpatelnya/rpi-home/backup/"+f.Name())
				if err != nil {
					sc.Logger.Errorf("Ошибка при попытке скопировать видео %s: %s", f.Name(), err.Error())

					return model.StatusNotSent, err
				}
				err = sc.Notifier.SendFile(fp, msg)
				if err != nil {
					sc.Logger.Errorf("Ошибка при попытке отправить видео %s: %s", f.Name(), err.Error())

					return model.StatusNotSent, err
				}
				sc.Logger.Infof("файл %s был отправлен в телеграм", fp)
				err = os.Remove(fp)
				if err != nil {
					sc.Logger.Errorf("Ошибка при попытке удалить файл: %s: %s", f.Name(), err.Error())

					return model.StatusNotSent, err
				}
			} else {
				err := sc.Notifier.SendText("Файл слишком велик чтобы его пересылать в Telegram. Вы можете его посмотреть через веб-интерфейс. Имя файла: " + f.Name())
				sentryhelper.Handle(sc.Logger, err, "Не удалось отправить текстовое сообщение о превышении размера видеофайла.")
				e.Status = model.StatusCanceled
				if err := sc.DB.Mongo.C("events").UpdateId(e.ID, e); err != nil {
					sc.Logger.Warningf("Ошибка сохранения события со статусом 'отменен'")
				}
			}
		}
	}

	return model.StatusSent, nil
}
