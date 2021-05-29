package sqlite3_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/model"
	"testing"
	"time"
)

func TestEventDataSQLite3_Save(t *testing.T) {
	startName := "test event"
	startDevice := "test device"
	newName := "new"
	e := &model.Event{
		Type:    model.TypeMotion,
		Status:  model.StatusNew,
		Name:    startName,
		Device:  startDevice,
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}

	assert.Nil(t, testContainer.Repo.Save(e))
	events := getAllEvents(t, testContainer.DB.SQLite3)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, startName, events[0].Name)
	assert.Equal(t, startDevice, events[0].Device)

	e.SqlID = 1
	e.Name = newName
	assert.Nil(t, testContainer.Repo.Save(e))

	events = getAllEvents(t, testContainer.DB.SQLite3)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, newName, events[0].Name)
	assert.Equal(t, startDevice, events[0].Device)
}
