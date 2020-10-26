package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"gopkg.in/mgo.v2/bson"
)

func TestAssertCreateMongoConnection(t *testing.T) {
	container, err := testhelpers.GetTestContainer(rootPath + configPath)
	assert.Nil(t, err)
	connectionSettings := container.AppConfig.Databases.MongoConnectionSettings

	t.Run("right config, wrong table", func(t *testing.T) {
		mongoConnection := config.AssertCreateMongoConnection(connectionSettings)
		var empty interface{}
		err = mongoConnection.C("random").Find(bson.M{}).One(&empty)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
