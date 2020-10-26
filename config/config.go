package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	EnvironmentDefault     = "default"
	EnvironmentProduction  = "production"
	EnvironmentTest        = "test"
	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"
)

var AppLevels = []string{EnvironmentDefault, EnvironmentProduction, EnvironmentTest, EnvironmentDevelopment, EnvironmentLocal}

// Config ...
type Config struct {
	Databases DbSettingsStruct `json:"databases"`
	Logger    Logger           `json:"logger"`
}

// DbSettingsStruct ...
type DbSettingsStruct struct {
	MongoConnectionSettings MongoConnectionSettings `json:"mongo"`
}

// MongoConnectionSettings ...
type MongoConnectionSettings struct {
	URI string `json:"uri"`
	DB  string `json:"db"`
}

// MongoConnection ...
type MongoConnection struct {
	session *mgo.Session
	setting MongoConnectionSettings
}

// ConnectionContainer ...
type ConnectionContainer struct {
	Mongo *MongoConnection
}

type Logger struct {
	LogLevel   string `json:"level"`
	ShowCaller bool   `json:"show_caller"`
}

// New ...
func New(p string) (*Config, error) {
	c, err := loadSettingsFromFile(p)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// C ...
func (db *MongoConnection) C(name string) *mgo.Collection {
	return db.session.Clone().DB(db.setting.DB).C(name)
}

func loadSettingsFromFile(path string) (*Config, error) {
	settingsJSON, err := readSettingsFile(path)

	if err != nil {
		return nil, err
	}

	return parseSettingsData(settingsJSON)
}

func readSettingsFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseSettingsData(settingsJSON []byte) (*Config, error) {
	var settings Config

	err := json.Unmarshal(settingsJSON, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// AssertCreateConnectionContainer ...
func (c *Config) AssertCreateConnectionContainer() *ConnectionContainer {
	mongoConnection := AssertCreateMongoConnection(c.Databases.MongoConnectionSettings)

	return &ConnectionContainer{Mongo: mongoConnection}
}

// CreateMongoConnection ...
func CreateMongoConnection(settings MongoConnectionSettings) (*MongoConnection, error) {
	session, err := mgo.Dial(settings.URI)

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(time.Second * 10)

	go mongoPing(session.DB(settings.DB))

	return &MongoConnection{session, settings}, nil
}

// AssertCreateMongoConnection ...
func AssertCreateMongoConnection(settings MongoConnectionSettings) *MongoConnection {
	log.Println("Connecting to mongo..")

	connection, err := CreateMongoConnection(settings)

	if err != nil {
		log.Println("Mongo connection error:", err)
		os.Exit(1)
	}

	return connection
}

func mongoPing(mg *mgo.Database) {
	errNum := 0

	for {
		err := mg.Session.Ping()
		if err != nil {
			mg.Session.Refresh()

			errNum++
		}

		if errNum > 5 {
			log.Println("To match error on mongo refresh connect. Shutdown")
			os.Exit(10)
		}

		time.Sleep(time.Second * 5)
	}
}
