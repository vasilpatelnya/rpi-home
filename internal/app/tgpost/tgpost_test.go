package tgpost

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendText(t *testing.T) {
	t.Run("empty text", func(t *testing.T) {
		err := SendText("")
		assert.Error(t, err, "отсутствует текст сообщения")
	})
	t.Run("simple text", func(t *testing.T) {
		err := SendText("simple text")
		assert.Nil(t, err)
	})
}

func TestSendFile(t *testing.T) {
	t.Run("not existed path", func(t *testing.T) {
		err := SendFile("./notExist.txt", "")
		assert.Error(t, err, "такого файла не существует или указанный путь неверен")
	})
	t.Run("empty path", func(t *testing.T) {
		err := SendFile("", "")
		assert.Error(t, err, "не указан путь к файлу")
	})
	t.Run("all right", func(t *testing.T) {
		err := SendFile("./tgpost_test.go", "тест")
		assert.Nil(t, err)
	})
}
