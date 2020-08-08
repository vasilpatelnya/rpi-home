package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testCfgPath = "../../../configs/test.env"

func TestNew(t *testing.T) {
	c, err := New(testCfgPath)
	assert.Nil(t, err)
	assert.Equal(t, c.AppMode, AppTest)
	c, err = New("wrong")
	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func TestConvertEnvVarToInt(t *testing.T) {
	t.Run("Wrong value", func(t *testing.T) {
		err := os.Setenv("DB_NAME", "example")
		assert.Nil(t, err)
		i, err := ConvertEnvVarToInt("", "DB_NAME")
		assert.NotNil(t, err)
		assert.Equal(t, 0, i)
	})
	t.Run("Right value", func(t *testing.T) {
		err := os.Setenv("MAIN_TICKER_TIME", "100")
		assert.Nil(t, err)
		i, err := ConvertEnvVarToInt("", "MAIN_TICKER_TIME")
		assert.Nil(t, err)
		assert.Equal(t, 100, i)
	})
}
