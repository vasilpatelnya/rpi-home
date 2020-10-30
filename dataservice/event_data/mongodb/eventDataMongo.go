package mongodb

import (
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type EventDataMongo struct {
	EventsCollection *mgo.Collection
	Logger           *logrus.Logger
}

func (data *EventDataMongo) GetAllByStatus(s int) ([]model.Event, error) {
	var events []model.Event
	if err := data.EventsCollection.Find(bson.M{"status": s}).All(&events); err != nil {
		data.Logger.Errorf("Ошибка при получении событий по статусу '%d': %s", s, err.Error())

		return nil, err
	}

	return events, nil
}

func (data *EventDataMongo) Save(e *model.Event) error {
	_, err := data.EventsCollection.Upsert(bson.M{"_id": e.ID}, e)
	if err != nil {
		data.Logger.Errorf("ошибка сохранения неотправленного события: %s", err.Error())

		return err
	}

	return nil
}

func (data *EventDataMongo) SaveUpdated(e *model.Event, status int) error {
	e.Status = status
	e.Updated = time.Now().UnixNano()

	return data.Save(e)
}
