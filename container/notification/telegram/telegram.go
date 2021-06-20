package telegram

import (
	"errors"
	"github.com/onrik/micha"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"os"
)

const (
	// LayoutISO ...
	LayoutISO = "2006-01-02"

	// StatusSent ...
	StatusSent = 1
	// StatusNotSent ...
	StatusNotSent = -1
)

// DirName - директория для мониторинга новых файлов.
type DirName string

// TGNotifier главная структура приложения.
type TGNotifier struct {
	michaAPI *micha.Bot
	chatID   micha.ChatID
}

// New ...
func New(token, chatID string) (notification.Notifier, error) {
	bot, err := micha.NewBot(token)
	if err != nil {
		return nil, err
	}

	go bot.Start()

	return &TGNotifier{
		michaAPI: bot,
		chatID:   micha.ChatID(chatID),
	}, nil
}

//SendText ...
func (tg *TGNotifier) SendText(t string) error {
	if len(t) == 0 {
		return errors.New("отсутствует текст сообщения")
	}

	_, err := tg.michaAPI.SendMessage(tg.chatID, t, nil)

	return err
}

//SendFile ...
func (tg *TGNotifier) SendFile(fp string, m string) error {
	if len(fp) == 0 {
		return errors.New("не указан путь к файлу")
	}
	exist := fs.Exists(fp)
	if !exist {
		return errors.New("такого файла не существует или указанный путь неверен")
	}
	file, err := os.Open(fp)
	if err != nil {
		return errors.New("reading file: error: " + err.Error())
	}
	_, err = tg.michaAPI.SendVideoFile(tg.chatID, file, fp, &micha.SendVideoOptions{Caption: m})

	return err
}
