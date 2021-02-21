package model

import (
	"fmt"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const (
	StatusNew   int = 0
	StatusReady int = 1
	StatusFail  int = -1

	TypeUndefined  int = 0
	TypeMotion     int = 1
	TypeMovieReady int = 2
)

type Event struct {
	ID      bson.ObjectId `bson:"_id"`
	Type    int           `bson:"type"`
	Status  int           `bson:"status"`
	Name    string        `bson:"name"`
	Device  string        `bson:"device"`
	Created int64         `bson:"created"` // timestamp Unix Nano
	Updated int64         `bson:"updated"` // timestamp Unix Nano
}

func New() *Event {
	return &Event{
		ID:      bson.NewObjectId(),
		Type:    TypeUndefined,
		Status:  StatusNew,
		Name:    "",
		Device:  "",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
}

func (e *Event) GetMotionMessage() string {
	text := fmt.Sprintf("Устройство %s зафиксировало движение! ", e.Device)
	var add string
	if e.Status == StatusFail {
		add = "Это сообщение отправлено с задержкой. Точное время срабатывания: " + time.Unix(e.Updated, 0).Format(time.Stamp)
	}

	return text + add
}

func (e *Event) GetVideoReadyMessage() string {
	tm := time.Unix(e.Created/1000000000, 0)
	msg := "Видео от " + tm.Format("2 January 2006 15:04")
	for en, ru := range translate.Months {
		msg = strings.Replace(msg, en, ru, -1)
	}

	return msg
}
