package store

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"testing"
)

func TestNew(t *testing.T) {
	path := "./../../../configs/" + config.TestCfgFilename
	c := config.New(path)
	t.Run("All right", func(t *testing.T) {
		s, err := New(c)
		assert.Nil(t, err)
		defer s.Connection.Close()
		assert.Nil(t, s.Connection.Ping())
	})
	t.Run("Bad connection URL", func(t *testing.T) {
		c.DbConnectionUrl = "mongodb://127.0.0.1:27018"
		s, err := New(c)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no reachable servers")
		assert.Nil(t, s)
	})
}
