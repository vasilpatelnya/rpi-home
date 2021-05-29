package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/mgo.v2"
)

const (
	EnvironmentDefault     = "default"
	EnvironmentProduction  = "production"
	EnvironmentTest        = "test"
	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"

	AppSettingsEnvName = "ENVIRONMENT"
)

var AppLevels = []string{EnvironmentDefault, EnvironmentProduction, EnvironmentTest, EnvironmentDevelopment, EnvironmentLocal}

// Config ...
type Config struct {
	Motion         MotionSettings   `json:"motion"`
	Periods        Periods          `json:"periods"`
	Databases      DbSettingsStruct `json:"databases"`
	Logger         Logger           `json:"logger"`
	SentrySettings SentrySettings   `json:"sentry_settings"`
	Notifier       Notifier         `json:"notifier"`
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
	MongoConnectionSettings   MongoConnectionSettings   `json:"mongo"`
	SQLite3ConnectionSettings SQLite3ConnectionSettings `json:"sqlite3"`
}

// MongoConnectionSettings ...
type MongoConnectionSettings struct {
	URI                 string        `json:"uri"`
	DB                  string        `json:"db"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
}

// SQLite3ConnectionSettings ...
type SQLite3ConnectionSettings struct {
	DBPath              string        `json:"db_path"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
}

// MongoConnection ...
type MongoConnection struct {
	session *mgo.Session
	setting MongoConnectionSettings
}

// SQLite3Connection ...
type SQLite3Connection struct {
	db         *sql.DB
	connection *sql.Conn
}

// ConnectionContainer ...
type ConnectionContainer struct {
	Mongo   *MongoConnection
	SQLite3 *SQLite3Connection
}

type Logger struct {
	LogLevel   string `json:"level"`
	ShowCaller bool   `json:"show_caller"`
}

type SentrySettings struct {
	SentryUrl string `json:"sentry_url"`
}

type Notifier struct {
	Type    string          `json:"type"`
	Options NotifierOptions `json:"options"`
}

type NotifierOptions struct {
	ChatID string `json:"chat_id"`
	Token  string `json:"token"`
}

// New ...
func New(p string) (*Config, error) {
	c, err := loadSettingsFromFile(p)
	if err != nil {
		return nil, err
	}

	log.Printf("Config info: %+v", c)

	return c, nil
}

func ParseEnvMode() (string, error) {
	env := os.Getenv(AppSettingsEnvName)
	if env == "" {
		return EnvironmentDefault, nil
	}
	match := false
	for _, level := range AppLevels {
		if env == level {
			match = true
			break
		}
	}
	if !match {
		msg := fmt.Sprintf("the specified operating mode (%s) of the application is incorrect, use the "+
			"following operating mode options: %s. Each mode of operation must correspond to the config of the same "+
			"name in the configs directory", env, strings.Join(AppLevels, ", "))

		return "", errors.New(msg)
	}

	return env, nil
}

// C ...
func (db *MongoConnection) C(name string) *mgo.Collection {
	return db.session.Clone().DB(db.setting.DB).C(name)
}

// C ...
func (db *SQLite3Connection) C() (*sql.Conn, *sql.DB) {
	return db.connection, db.db
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
	sqlite3Connection := AssertCreateSQLite3Connection(c.Databases.SQLite3ConnectionSettings)

	return &ConnectionContainer{Mongo: mongoConnection, SQLite3: sqlite3Connection}
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

// CreateSQLite3Connection ...
func CreateSQLite3Connection(settings SQLite3ConnectionSettings) (*SQLite3Connection, error) {
	rootPath, err := fs.RootPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s", rootPath, settings.DBPath))

	if err != nil {
		return nil, err
	}

	connection, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	go sqlite3Ping(db, settings)

	return &SQLite3Connection{connection: connection, db: db}, nil
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
			log.Fatal("Превышено количество попыток подключения к Mongo DB. Завершение работы.")
		}

		time.Sleep(time.Second * settings.TimeBetweenAttempts)
	}
}

func sqlite3Ping(sqlite3 *sql.DB, settings SQLite3ConnectionSettings) {
	errNum := 0

	for {
		err := sqlite3.Ping()
		if err != nil {
			errNum++
		}

		if errNum > settings.ConnectAttempts {
			log.Fatal("Превышено количество попыток подключения к SQLite3. Завершение работы.")
		}

		time.Sleep(time.Second * settings.TimeBetweenAttempts)
	}
}

// AssertCreateSQLite3Connection ...
func AssertCreateSQLite3Connection(settings SQLite3ConnectionSettings) *SQLite3Connection {
	log.Println("Устанавливаем соединение с SQLite 3...")

	connection, err := CreateSQLite3Connection(settings)

	if err != nil {
		log.Println("Ошибка при создании подключения к БД.", err)
		os.Exit(1)
	}

	return connection
}
