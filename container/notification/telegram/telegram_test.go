package telegram

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := os.Setenv("APP_MODE", "test")
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func TestSendText(t *testing.T) {
	tg := TGNotifier{}
	t.Run("empty text", func(t *testing.T) {
		err := tg.SendText("")
		assert.Error(t, err, "отсутствует текст сообщения")
	})
	t.Run("simple text", func(t *testing.T) {
		err := tg.SendText("simple text")
		assert.Nil(t, err)
	})
}

func TestSendFile(t *testing.T) {
	tg := TGNotifier{}
	t.Run("not existed path", func(t *testing.T) {
		err := tg.SendFile("./notExist.txt", "")
		assert.Error(t, err, "такого файла не существует или указанный путь неверен")
	})
	t.Run("empty path", func(t *testing.T) {
		err := tg.SendFile("", "")
		assert.Error(t, err, "не указан путь к файлу")
	})
	t.Run("all right", func(t *testing.T) {
		err := tg.SendFile("./telegram_test.go", "тест")
		assert.Nil(t, err)
	})
}
