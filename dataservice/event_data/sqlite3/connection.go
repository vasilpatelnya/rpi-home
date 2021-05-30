package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/tool/fs"
	"log"
	"os"
	"time"
)

// SQLite3Connection ...
type SQLite3Connection struct {
	db         *sql.DB
	connection *sql.Conn
}

// AssertCreateSQLite3Connection ...
func AssertCreateSQLite3Connection(settings *config.SQLite3ConnectionSettings) *SQLite3Connection {
	log.Println("Устанавливаем соединение с SQLite 3...")

	connection, err := CreateSQLite3Connection(settings)

	if err != nil {
		log.Println("Ошибка при создании подключения к БД.", err)
		os.Exit(1)
	}

	return connection
}

// C ...
func (db *SQLite3Connection) C() (*sql.Conn, *sql.DB) {
	return db.connection, db.db
}

// CreateSQLite3Connection ...
func CreateSQLite3Connection(settings *config.SQLite3ConnectionSettings) (*SQLite3Connection, error) {
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

func sqlite3Ping(sqlite3 *sql.DB, settings *config.SQLite3ConnectionSettings) {
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
