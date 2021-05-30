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

var AppLevels = []string{EnvironmentDefault, EnvironmentProduction, EnvironmentTest, EnvironmentDevelopment, EnvironmentLocal}

// Config ...
type Config struct {
	Motion         MotionSettings   `json:"motion"`
	Periods        Periods          `json:"periods"`
	Database       DbSettingsStruct `json:"database"`
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
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
}

// MongoSettings ...
type MongoSettings struct {
	Type     string                  `json:"type"`
	Settings MongoConnectionSettings `json:"settings"`
}

type MongoConnectionSettings struct {
	URI                 string        `json:"uri"`
	DB                  string        `json:"db"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
}

// SQLite3Settings ...
type SQLite3Settings struct {
	Type     string                    `json:"type"`
	Settings SQLite3ConnectionSettings `json:"settings"`
}

// SQLite3ConnectionSettings ...
type SQLite3ConnectionSettings struct {
	DBPath              string        `json:"db_path"`
	ConnectAttempts     int           `json:"connect_attempts"`
	TimeBetweenAttempts time.Duration `json:"time_between_attempts"`
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
