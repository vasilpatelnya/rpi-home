package rpi_detector_mongo

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/config"
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/store"
	"gopkg.in/mgo.v2/bson"
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
	var events []Event
	c := config.New("./../../../configs/test.env")
	s, err := store.New(c)
	assert.Nil(t, err)
	err = s.Collection.Find(bson.M{"status": StatusReady}).All(&events)
	assert.True(t, len(events) > 0)
}
