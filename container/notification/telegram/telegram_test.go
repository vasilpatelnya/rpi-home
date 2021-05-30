package telegram_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/notification/telegram"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	var err error
	env, err := config.ParseEnvMode()
	if err != nil {
		log.Fatal(err)
	}
	rootPath, err := fs.RootPath()
	if err != nil {
		log.Fatal(err)
	}
	testConfig, err = config.New(fmt.Sprintf("%s/config/%s.json", rootPath, env))
	if err != nil {
		log.Fatal(err)
	}
	if !testConfig.Notifier.IsUsing {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestSendText(t *testing.T) {
	tg, err := telegram.New(testConfig.Notifier.Options.Token, testConfig.Notifier.Options.ChatID)
	assert.Nil(t, err)
	assert.Nil(t, tg.SendText("test message from Tosya"))
}

func TestSendFile(t *testing.T) {
	tg, err := telegram.New(testConfig.Notifier.Options.Token, testConfig.Notifier.Options.ChatID)
	assert.Nil(t, err)

	src := testhelpers.GetTestDataDir() + "/2020-10-31/test.mp4"
	err = ioutil.WriteFile(src, []byte("test data"), 0777)
	assert.Nil(t, err)

	t.Log(src)
	assert.Nil(t, tg.SendFile(src, "test message"))
}
