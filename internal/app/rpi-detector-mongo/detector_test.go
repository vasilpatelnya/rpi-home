package rpi_detector_mongo

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/config"
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/store"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	e := New()
	assert.True(t, reflect.DeepEqual(e, &Event{
		Type:    TypeUndefined,
		Status:  StatusNew,
		Name:    "",
		Device:  "",
		Created: e.Created,
		Updated: e.Updated,
	}))
}

func TestEvent_GetAllByStatus(t *testing.T) {
	c := config.New("./../../../configs/test.env")
	s, err := store.New(c)
	assert.Nil(t, err)
	_, err = GetAllByStatus(s.Collection, StatusNew)
	assert.Nil(t, err)
}
