package dataservice

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/mongodb"
	"github.com/vasilpatelnya/rpi-home/dataservice/event_data/sqlite3"
)

// ConnectionContainer ...
type ConnectionContainer struct {
	Mongo   *mongodb.MongoConnection
	SQLite3 *sqlite3.SQLite3Connection
}

// NewConnectionContainer ...
func NewConnectionContainer(c config.DbSettings) (*ConnectionContainer, error) {
	var mongoConnection *mongodb.MongoConnection
	var sqlite3Connection *sqlite3.SQLite3Connection

	// todo refactor this shit!!!
	if c.Type == config.DbTypeMongo {
		raw, ok := c.Settings.(map[string]interface{})
		if !ok {
			return nil, errors.New("(mongo) wrong settings type")
		}
		uri, okURI := raw["uri"].(string)
		db, okDB := raw["db"].(string)
		connectAttempts, okCA := raw["connect_attempts"].(int)
		timeBetweenAttempts, okTBA := raw["time_between_attempts"].(time.Duration)

		if !okURI || !okDB {
			msg := fmt.Sprintf("uri: %v, db: %v, ca: %v, tba: %v", okURI, okDB, okCA, okTBA)
			return nil, errors.New("parse settings error. " + msg)
		}

		connectAttempts = 10
		timeBetweenAttempts = time.Second * 20

		settings := config.MongoSettings{
			URI:                 uri,
			DB:                  db,
			ConnectAttempts:     connectAttempts,
			TimeBetweenAttempts: timeBetweenAttempts,
		}

		mongoConnection = mongodb.AssertCreateMongoConnection(&settings)
	}
	if c.Type == config.DbTypeSQLite3 {
		raw, ok := c.Settings.(map[string]interface{})
		if !ok {
			return nil, errors.New("(sqlite3) wrong settings type")
		}
		dbPath, okURI := raw["db_path"].(string)

		if !okURI {
			msg := fmt.Sprintf("uri: %v", okURI)
			return nil, errors.New("parse settings error. " + msg)
		}

		settings := config.SQLite3Settings{
			DBPath:              dbPath,
			ConnectAttempts:     10,
			TimeBetweenAttempts: time.Second * 20,
		}

		sqlite3Connection = sqlite3.AssertCreateSQLite3Connection(&settings)
	}

	return &ConnectionContainer{Mongo: mongoConnection, SQLite3: sqlite3Connection}, nil
}
