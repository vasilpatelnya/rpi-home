package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

// Config ...
type Config struct {
	Motion         MotionSettings   `json:"motion"`
	Periods        Periods          `json:"periods"`
	Databases      DbSettingsStruct `json:"databases"`
	Logger         Logger           `json:"logger"`
	SentrySettings SentrySettings   `json:"sentry_settings"`
}

type MotionSettings struct {
	MoviesDirCam1 string `json:"movies_dir_cam_1"`
	FileExtension string `json:"file_extension"`
}

type Periods struct {
	MainTickerTime time.Duration `json:"main_ticker_time"`
}

// DbSettingsStruct ...
type DbSettingsStruct struct {
	MongoConnectionSettings MongoConnectionSettings `json:"mongo"`
}

// MongoConnectionSettings ...
type MongoConnectionSettings struct {
	URI                 string        `json:"uri"`
	DB                  string        `json:"db"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
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

type SentrySettings struct {
	SentryUrl string `json:"sentry_url"`
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

	go mongoPing(session.DB(settings.DB), settings)

	return &MongoConnection{session, settings}, nil
}

// AssertCreateMongoConnection ...
func AssertCreateMongoConnection(settings MongoConnectionSettings) *MongoConnection {
	log.Println("Устанавливаем соединение с Mongo DB...")

	connection, err := CreateMongoConnection(settings)

	if err != nil {
		log.Println("Ошибка при создании подключения к БД.", err)
		os.Exit(1)
	}

	return connection
}

func mongoPing(mg *mgo.Database, settings MongoConnectionSettings) {
	errNum := 0

	for {
		err := mg.Session.Ping()
		if err != nil {
			mg.Session.Refresh()

			errNum++
		}

		if errNum > settings.ConnectAttempts {
			log.Println("Превышено количество попыток подключения к Mongo DB. Завершение работы.")
			os.Exit(10)
		}

		time.Sleep(time.Second * settings.TimeBetweenAttempts)
	}
}
