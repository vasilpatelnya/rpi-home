package mongodb

import (
	"github.com/vasilpatelnya/rpi-home/config"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

// MongoConnection ...
type MongoConnection struct {
	session *mgo.Session
	setting *config.MongoSettings
}

// CreateMongoConnection ...
func CreateMongoConnection(c *config.MongoSettings) (*MongoConnection, error) {
	session, err := mgo.Dial(c.Settings.URI)

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(time.Second * 10)

	go mongoPing(session.DB(c.Settings.DB), c)

	return &MongoConnection{session, c}, nil
}

// AssertCreateMongoConnection ...
func AssertCreateMongoConnection(settings *config.MongoSettings) *MongoConnection {
	log.Println("Устанавливаем соединение с Mongo DB...")

	connection, err := CreateMongoConnection(settings)

	if err != nil {
		log.Println("Ошибка при создании подключения к БД.", err)
		os.Exit(1)
	}

	return connection
}

func mongoPing(mg *mgo.Database, c *config.MongoSettings) {
	errNum := 0

	for {
		err := mg.Session.Ping()
		if err != nil {
			mg.Session.Refresh()

			errNum++
		}

		if errNum > c.Settings.ConnectAttempts {
			log.Fatal("Превышено количество попыток подключения к Mongo DB. Завершение работы.")
		}

		time.Sleep(time.Second * c.Settings.TimeBetweenAttempts)
	}
}

// C ...
func (db *MongoConnection) C(name string) *mgo.Collection {
	return db.session.Clone().DB(db.setting.Settings.DB).C(name)
}
