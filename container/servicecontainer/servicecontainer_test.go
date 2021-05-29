package servicecontainer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
)

func TestServiceContainer_InitApp(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		c := getTestServiceContainer()
		err := c.InitApp()
		assert.Nil(t, err)
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
		assert.Equal(t, "not a valid logrus Level: \"\"", app.InitLogger().Error())
		assert.Equal(t, "info", app.Logger.Level.String())
		assert.Equal(t, c.Logger.ShowCaller, app.Logger.ReportCaller)
	})
}
