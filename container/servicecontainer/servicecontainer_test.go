package servicecontainer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotNil(t, app.Logger)
}
