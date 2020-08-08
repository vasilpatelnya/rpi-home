package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const (
	AppDev          = "dev"
	AppProd         = "prod"
	AppTest         = "test"
	TestCfgFilename = "test.env"
)

type Config struct {
	AppMode               string
	MainTickerTime        uint16
	DbConnectionUrl       string
	DbName                string
	DbTable               string
	MoviesDirCamera1      string
	FileExtension         string
	DbConnectAttempts     int
	DbTimeBetweenAttempts int
	SentryUrl             string
}

func New(p string) *Config {
	if err := godotenv.Load(p); err != nil {
		log.Fatal("Ошибка загрузки env файла.", p)
	}
	mtt := ConvertEnvVarToInt(p, "MAIN_TICKER_TIME")
	ca := ConvertEnvVarToInt(p, "DB_CONNECT_ATTEMPTS")
	tba := ConvertEnvVarToInt(p, "DB_TIME_BETWEEN_ATTEMPTS")

	return &Config{
		AppMode:               os.Getenv("APP_MODE"),
		MainTickerTime:        uint16(mtt),
		DbConnectionUrl:       os.Getenv("DB_CONNECTION_URL"),
		DbName:                os.Getenv("DB_NAME"),
		DbTable:               os.Getenv("DB_TABLE"),
		MoviesDirCamera1:      os.Getenv("MOVIES_DIR_CAMERA1"),
		FileExtension:         os.Getenv("FILE_EXTENSION"),
		DbConnectAttempts:     ca,
		DbTimeBetweenAttempts: tba,
		SentryUrl:             os.Getenv("SENTRY_URL"),
	}
}

func ConvertEnvVarToInt(p string, s string) int {
	res, err := strconv.Atoi(os.Getenv(s))
	if err != nil {
		log.Fatalf("Ошибка чтения параметра %s в конфигурационном файле: %s", s, p)
	}

	return res
}
