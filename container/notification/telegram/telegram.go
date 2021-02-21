package telegram

import (
	"errors"
	"fmt"
	"github.com/onrik/micha"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"log"
)

const (
	appPath = "/usr/local/bin/telegram-send"

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
func New(token, chatID string) notification.Notifier {
	bot, err := micha.NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	go bot.Start()

	return &TGNotifier{
		michaAPI: bot,
		chatID:   micha.ChatID(chatID),
	}
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
	caption := ""
	if len(m) > 0 {
		caption = fmt.Sprintf(`--caption "%s"`, m)
	}
	_, err := tg.michaAPI.SendVideo(tg.chatID, fp, &micha.SendVideoOptions{Caption: caption})

	return err
}
