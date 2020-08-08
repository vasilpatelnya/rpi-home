package tgpost

import (
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	appPath   = "/usr/local/bin/telegram-send"
	LayoutISO = "2006-01-02"

	StatusSent    = 1
	StatusNotSent = -1
)

// DirName - директория для мониторинга новых файлов.
type DirName string

// TgPost главная структура приложения.
type TgPost struct {
	Message       string
	Filepath      string
	DirName       string
	FileExtension string
}

//SendText ...
func SendText(t string) error {
	if len(t) == 0 {
		return errors.New("отсутствует текст сообщения")
	}
	if os.Getenv("APP_MODE") == config.AppTest {
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
func SendFile(fp string, m string) error {
	if len(fp) == 0 {
		return errors.New("не указан путь к файлу")
	}
	exist := exists(fp)
	if !exist {
		return errors.New("такого файла не существует или указанный путь неверен")
	}
	if os.Getenv("APP_MODE") == config.AppTest {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

//GetDirList ...
func GetTodayFileList(dirname string) ([]os.FileInfo, error) {
	if exists(dirname) {
		files, err := ioutil.ReadDir(GetTodayPath(dirname))
		if err != nil {
			return nil, err
		}

		return files, nil
	}

	return nil, errors.New(dirname + " directory is not exist")
}

// GetTodayDir ...
func GetTodayDir() string {
	t := time.Now()

	return t.Format(LayoutISO)
}

// GetTodayPath ...
func GetTodayPath(dirname string) string {
	return fmt.Sprintf("%s/%s", dirname, GetTodayDir())
}
