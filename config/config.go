package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var AppLevels = []string{
	EnvironmentDefault, EnvironmentProduction, EnvironmentTest,
	EnvironmentDevelopment, EnvironmentLocal, EnvironmentCiMongo,
	EnvironmentCiSQLites3, EnvironmentDocker,
}

// Config ...
type Config struct {
	Motion         MotionSettings `json:"motion"`
	Periods        Periods        `json:"periods"`
	Database       DbSettings     `json:"database"`
	ApiServer      ApiSettings    `json:"api_server"`
	Logger         Logger         `json:"logger"`
	Notifier       Notifier       `json:"notifier"`
	SentrySettings SentrySettings `json:"sentry_settings"`
}

type MotionSettings struct {
	MoviesDirCam1 string `json:"movies_dir_cam_1"`
	FileExtension string `json:"file_extension"`
}

type Periods struct {
	MainTickerTime time.Duration `json:"main_ticker_time"`
}

type DbSettings struct {
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
}

// MongoSettings ...
type MongoSettings struct {
	URI                 string        `json:"uri"`
	DB                  string        `json:"db"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
}

// SQLite3Settings ...
type SQLite3Settings struct {
	DBPath              string        `json:"db_path"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
}

type ApiSettings struct {
	Port   int    `json:"port"`
	ApiKey string `json:"api_key"`
}

type Logger struct {
	LogLevel   string `json:"level"`
	ShowCaller bool   `json:"show_caller"`
}

type SentrySettings struct {
	SentryUrl string `json:"sentry_url"`
}

type Notifier struct {
	IsUsing bool            `json:"is_using"`
	Type    string          `json:"type"`
	Options NotifierOptions `json:"options"`
}

type NotifierOptions struct {
	ChatID string `json:"chat_id"`
	Token  string `json:"token"`
}

// New ...
func New(p string) (*Config, error) {
	log.Printf("Try to load settings from file [%s]", p)
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
