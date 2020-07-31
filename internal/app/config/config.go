package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const TestCfgFilename = "test.env"

type Config struct {
	MainTickerTime   uint16
	DbConnectionUrl  string
	DbName           string
	DbTable          string
	MoviesDirCamera1 string
	FileExtension    string
}

func New(p string) *Config {
	if err := godotenv.Load(p); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	mtt, err := strconv.Atoi(os.Getenv("MAIN_TICKER_TIME"))
	if err != nil {
		log.Fatal("Ошибка чтения параметра MAIN_TICKER_TIME в конфигурационном файле:", p)
	}

	return &Config{
		MainTickerTime:   uint16(mtt),
		DbConnectionUrl:  os.Getenv("DB_CONNECTION_URL"),
		DbName:           os.Getenv("DB_NAME"),
		DbTable:          os.Getenv("DB_TABLE"),
		MoviesDirCamera1: os.Getenv("MOVIES_DIR_CAMERA1"),
		FileExtension:    os.Getenv("FILE_EXTENSION"),
	}
}
