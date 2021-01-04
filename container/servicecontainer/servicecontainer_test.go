package servicecontainer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
)

func TestServiceContainer_InitApp(t *testing.T) {
	t.Run("Right config path", func(t *testing.T) {
		c := getTestServiceContainer()
		err := c.InitApp(RightConfigPath)
		assert.Nil(t, err)
	})
	t.Run("Wrong config path", func(t *testing.T) {
		c := getTestServiceContainer()
		err := c.InitApp(WrongConfigPath)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})
}

func TestServiceContainer_InitLogger(t *testing.T) {
	app := getTestServiceContainer()

	c, err := config.New(RightConfigPath)
	assert.Nil(t, err)
	app.AppConfig = c

	t.Run("right values", func(t *testing.T) {
		assert.Nil(t, app.InitLogger())
		assert.Equal(t, c.Logger.LogLevel, app.Logger.Level.String())
		assert.Equal(t, c.Logger.ShowCaller, app.Logger.ReportCaller)
	})

	t.Run("empty log level", func(t *testing.T) {
		app.AppConfig.Logger.LogLevel = ""
		assert.Nil(t, app.InitLogger())
		assert.Equal(t, "info", app.Logger.Level.String())
		assert.Equal(t, c.Logger.ShowCaller, app.Logger.ReportCaller)
	})
}
