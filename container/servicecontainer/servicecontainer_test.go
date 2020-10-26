package servicecontainer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
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

func getTestServiceContainer() servicecontainer.ServiceContainer {
	return servicecontainer.ServiceContainer{AppConfig: &config.Config{}}
}
