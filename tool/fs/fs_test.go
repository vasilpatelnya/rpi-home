package fs_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"io/ioutil"
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

func TestCopyFile(t *testing.T) {
	newFilePath := testhelpers.GetTestDataDir() + "/test.mp4"
	err := ioutil.WriteFile(newFilePath, []byte("test data"), 0777)
	assert.Nil(t, err)

	dst := testhelpers.GetTestDataDir() + "/test_copy.mp4"
	assert.Nil(t, fs.CopyFile(newFilePath, dst))
}
