package config_test

import (
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"gopkg.in/mgo.v2/bson"
)

func TestAssertCreateMongoConnection(t *testing.T) {
	container, err := testhelpers.GetTestContainer()
	assert.Nil(t, err)
	connectionSettings, isMongoSettings := container.AppConfig.Database.Settings.(config.MongoSettings)

	if !isMongoSettings {
		t.Skip("config not contain mongo settings")
	}

	t.Run("right config, wrong table", func(t *testing.T) {
		mongoConnection := mongodb.AssertCreateMongoConnection(&connectionSettings)
		var empty interface{}
		err = mongoConnection.C("random").Find(bson.M{}).One(&empty)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
