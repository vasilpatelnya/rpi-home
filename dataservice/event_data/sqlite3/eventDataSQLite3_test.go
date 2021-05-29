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

func TestEventDataSQLite3_GetAllByStatus(t *testing.T) {
	clearTable(t, testContainer.DB.SQLite3, "events")
	event1 := &model.Event{
		Type:    model.TypeMotion,
		Status:  model.StatusNew,
		Name:    "name 1",
		Device:  "device 1",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	event2 := &model.Event{
		Type:    model.TypeMotion,
		Status:  model.StatusReady,
		Name:    "name 2",
		Device:  "device 2",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	event3 := &model.Event{
		Type:    model.TypeMotion,
		Status:  model.StatusNew,
		Name:    "name 3",
		Device:  "device 3",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	for _, event := range []*model.Event{event1, event2, event3} {
		assert.Nil(t, testContainer.Repo.Save(event))
	}

	events, err := testContainer.Repo.GetAllByStatus(model.StatusNew)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(events))
}
