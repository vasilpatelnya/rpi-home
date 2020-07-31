package store

import (
	"gitlab.com/vasilpatelnya/rpi-home/internal/app/config"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

const (
	StatusUpdated    = 11
	StatusNotUpdated = -11
)

type Store struct {
	Connection *mgo.Session
	Collection *mgo.Collection
}

func New(c *config.Config) (*Store, error) {
	var client *mgo.Session
	for {
		var err error
		client, err = getConnection(c)

		if err != nil {
			log.Println("Ошибка подключения к базе")
			time.Sleep(time.Second * 20)
			continue
		}

		break
	}

	log.Println("Подключено к MongoDB!")
	collection := client.DB(c.DbName).C(c.DbTable)

	return &Store{
		Connection: client,
		Collection: collection,
	}, nil
}

func getConnection(c *config.Config) (*mgo.Session, error) {
	client, err := mgo.Dial(c.DbConnectionUrl)
	if err != nil {
		return nil, err
	}

	return client, nil
}
