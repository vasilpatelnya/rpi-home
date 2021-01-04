package fs_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"testing"
	"time"
)

func TestGetTodayDir(t *testing.T) {
	now := time.Now()
	assert.Equal(t, fs.GetTodayDir(telegram.LayoutISO), now.Format(telegram.LayoutISO))
}

func TestGetTodayPath(t *testing.T) {
	now := time.Now().Format(telegram.LayoutISO)
	assert.Equal(t, fs.GetTodayPath("test", telegram.LayoutISO), "test/"+now)
}

func TestGetTodayFileList(t *testing.T) {
	dir := "wrong"
	_, err := fs.GetTodayFileList(dir, telegram.LayoutISO)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), dir+" directory is not exist")
}
