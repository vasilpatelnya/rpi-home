package usecase

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	sentryhelper "github.com/vasilpatelnya/rpi-home/container/sentry-helper"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/model"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"os"
	"path/filepath"
)

const (
	MaxSize          int64 = 50 * 1024 * 1024 // 50mb
	StatusUpdated          = 11
	StatusNotUpdated       = -11
)

type EventHandleOpts struct {
	TargetDir string // todo переделать на аргумент с несколькими путями к записям: например для нескольких камер
	BackupDir string
	Ext       string
	Repo      dataservice.EventData
	Notifier  notification.Notifier
	Logger    *logrus.Logger
}

func EventHandle(opts EventHandleOpts) {
	handleEvents(model.StatusFail, opts)
	handleEvents(model.StatusNew, opts)
}

func handleEvents(status int, opts EventHandleOpts) {
	events, err := opts.Repo.GetAllByStatus(status)
	if err != nil {
		sentryhelper.Handle(opts.Logger, err, "Ошибка получения записей событий из БД")
	}

	if len(events) > 0 {
		for _, e := range events {
			switch e.Type {
			case model.TypeMotion:
				status, err := handleMotionAlarm(opts, &e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(opts.Logger, err, msg)
					if status == model.StatusNotSent {
						e.Status = model.StatusFail
						err = opts.Repo.Save(&e)
						if err != nil {
							msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
							sentryhelper.Handle(opts.Logger, err, msg)
						}
						continue
					}
				}
				e.Status = model.StatusReady
				err = opts.Repo.Save(&e)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(opts.Logger, err, msg)
				}
			case model.TypeMovieReady:
				opts.Logger.Info("Видео готово!")
				e.Status, err = handleMotionReady(&e, opts)
				if err != nil {
					msg := fmt.Sprintf("Ошибка обработки события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(opts.Logger, err, msg)
				}

				err = opts.Repo.SaveUpdated(&e, e.Status)
				if err != nil {
					msg := fmt.Sprintf("Ошибка сохранения события: %s %s", e.Name, err.Error())
					sentryhelper.Handle(opts.Logger, err, msg)
				}
			}
		}
	}
}

func handleMotionAlarm(opts EventHandleOpts, e *model.Event) (int, error) {
	err := opts.Notifier.SendText(e.GetMotionMessage())
	if err != nil {
		return model.StatusNotSent, errors.New("ошибка отправки текста о срабатывании")
	}
	err = opts.Repo.Save(e)
	if err != nil {
		return StatusNotUpdated, errors.New(fmt.Sprintf("ошибка обновления записи в БД, id записи: %s", e.ID.Hex()))
	}

	return StatusUpdated, nil
}

func handleMotionReady(e *model.Event, opts EventHandleOpts) (int, error) {
	l, err := fs.GetTodayFileList(opts.TargetDir, model.LayoutISO)
	if err != nil {
		opts.Logger.Errorf("Ошибка получения списка файлов в директории %s: %s", opts.TargetDir, err.Error())

		return model.StatusNotSent, err
	}
	for _, f := range l {
		todayDir := fs.GetTodayDir(model.LayoutISO)
		fp := fmt.Sprintf("%s/%s/%s", opts.TargetDir, todayDir, f.Name())
		ext := filepath.Ext(f.Name())
		if ext == opts.Ext && f.Size() > 0 {
			if f.Size() < MaxSize {
				msg := e.GetVideoReadyMessage()
				dstPath := opts.BackupDir + "/" + f.Name()
				err = fs.CopyFile(fp, dstPath)
				if err != nil {
					opts.Logger.Errorf("Ошибка при попытке скопировать видео \nиз %s \nв %s\n Ошибка: %s", fp, dstPath, err.Error())

					return model.StatusNotSent, err
				}
				err = opts.Notifier.SendFile(fp, msg)
				if err != nil {
					opts.Logger.Errorf("Ошибка при попытке отправить видео %s: %s", f.Name(), err.Error())

					return model.StatusNotSent, err
				}
				opts.Logger.Infof("файл %s был отправлен в телеграм", fp)
				err = os.Remove(fp)
				if err != nil {
					opts.Logger.Errorf("Ошибка при попытке удалить файл: %s: %s", f.Name(), err.Error())

					return model.StatusNotSent, err
				}
			} else {
				err := opts.Notifier.SendText("Файл слишком велик чтобы его пересылать в Telegram. Вы можете его посмотреть через веб-интерфейс. Имя файла: " + f.Name())
				sentryhelper.Handle(opts.Logger, err, "Не удалось отправить текстовое сообщение о превышении размера видеофайла.")
				e.Status = model.StatusCanceled
				if err := opts.Repo.Save(e); err != nil {
					opts.Logger.Warningf("Ошибка сохранения события со статусом 'отменен'")
				}
			}
		}
	}

	return model.StatusSent, nil
}
