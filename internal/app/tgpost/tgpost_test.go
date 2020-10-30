package tgpost

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := os.Setenv("APP_MODE", "test")
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

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

func TestGetTodayDir(t *testing.T) {
	now := time.Now()
	assert.Equal(t, GetTodayDir(), now.Format(LayoutISO))
}

func TestGetTodayPath(t *testing.T) {
	now := time.Now().Format(LayoutISO)
	assert.Equal(t, GetTodayPath("test"), "test/"+now)
}

func TestGetTodayFileList(t *testing.T) {
	dir := "wrong"
	_, err := GetTodayFileList(dir)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), dir+" directory is not exist")
}
