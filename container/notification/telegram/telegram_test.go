package telegram_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	var err error
	testConfig, err = config.New("./../../../config/docker.json")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestSendText(t *testing.T) {
	tg := telegram.New(testConfig.Notifier.Options.Token, testConfig.Notifier.Options.ChatID)
	assert.Nil(t, tg.SendText("test message from Tosya"))
}

func TestSendFile(t *testing.T) {
	tg := telegram.New(testConfig.Notifier.Options.Token, testConfig.Notifier.Options.ChatID)

	src := testhelpers.GetTestDataDir() + "/2020-10-31/test.mp4"
	err := ioutil.WriteFile(src, []byte("test data"), 0777)
	assert.Nil(t, err)

	t.Log(src)
	assert.Nil(t, tg.SendFile(src, "test message"))
}
