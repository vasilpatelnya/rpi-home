package sqlite3_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/model"
	"testing"
	"time"
)

func TestEventDataSQLite3_Save(t *testing.T) {
	e := &model.Event{
		Type:    model.TypeMotion,
		Status:  model.StatusNew,
		Name:    "test event",
		Device:  "test device",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}

	assert.Nil(t, testContainer.Repo.Save(e))
	e.SqlID = 1
	e.Name = "new"
	assert.Nil(t, testContainer.Repo.Save(e))
}
