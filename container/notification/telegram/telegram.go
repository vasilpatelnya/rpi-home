package telegram

import (
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/container/notification"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"log"
	"os"
	"os/exec"
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
type TGNotifier struct{}

func New() notification.Notifier {
	return new(TGNotifier)
}

//SendText ...
func (tg *TGNotifier) SendText(t string) error {
	if len(t) == 0 {
		return errors.New("отсутствует текст сообщения")
	}
	if os.Getenv("APP_MODE") == "test" {
		log.Println("Вы находитесь в тестовом режиме. Отправка файлов игнорируется.")
		return nil
	}
	cmd := fmt.Sprintf(`%s "%s"`, appPath, t)
	if err := exec.Command("/bin/bash", "-c", cmd).Run(); err != nil {
		return err
	}

	return nil
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
	if os.Getenv("APP_MODE") == "test" {
		log.Println("Вы находитесь в тестовом режиме. Отправка файлов игнорируется.")
		return nil
	}
	caption := ""
	if len(m) > 0 {
		caption = fmt.Sprintf(`--caption "%s"`, m)
	}
	cmd := fmt.Sprintf(`%s --file %s %s`, appPath, fp, caption)
	if err := exec.Command("/bin/bash", "-c", cmd).Run(); err != nil {
		log.Println("Error after command '"+cmd+"':", err.Error())
		return err
	}

	return nil
}
