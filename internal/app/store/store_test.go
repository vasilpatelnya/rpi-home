package store

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/config"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	path := "./../../../configs/" + config.TestCfgFilename
	log.Println(path)
	c := config.New(path)
	s, err := New(c)
	assert.Nil(t, err)
	defer s.Connection.Close()
	assert.Nil(t, s.Connection.Ping())
}
