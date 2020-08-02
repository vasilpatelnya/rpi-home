package rpi_detector_mongo

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"github.com/vasilpatelnya/rpi-home/internal/app/store"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
	"time"
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

func TestEvent_Save(t *testing.T) {
	c := config.New("./../../../configs/test.env")
	s, err := store.New(c)
	assert.Nil(t, err)
	event := &Event{
		ID:      bson.NewObjectId(),
		Type:    1,
		Status:  1,
		Name:    "test",
		Device:  "test",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	err = event.Save(s.Collection)
	assert.Nil(t, err)
}
