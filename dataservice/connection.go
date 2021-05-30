package dataservice

import (
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/sqlite3"
)

// ConnectionContainer ...
type ConnectionContainer struct {
	Mongo   *mongodb.MongoConnection
	SQLite3 *sqlite3.SQLite3Connection
}

// AssertCreateConnectionContainer ...
func AssertCreateConnectionContainer(settings interface{}) *ConnectionContainer {
	var mongoConnection *mongodb.MongoConnection
	var sqlite3Connection *sqlite3.SQLite3Connection
	mongoSettings, isMongoSettings := settings.(*config.MongoConnectionSettings)
	sqlite3Settings, isSQLite3Settings := settings.(*config.SQLite3ConnectionSettings)
	if isMongoSettings && mongoSettings != nil {
		mongoConnection = mongodb.AssertCreateMongoConnection(mongoSettings)
	}
	if isSQLite3Settings && sqlite3Settings != nil {
		sqlite3Connection = sqlite3.AssertCreateSQLite3Connection(sqlite3Settings)
	}

	return &ConnectionContainer{Mongo: mongoConnection, SQLite3: sqlite3Connection}
}
