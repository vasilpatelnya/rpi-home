package daemon

import (
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	rpidetectormongo "github.com/vasilpatelnya/rpi-home/internal/app/rpi-detector-mongo"
	"github.com/vasilpatelnya/rpi-home/internal/app/store"
	"github.com/vasilpatelnya/rpi-home/internal/app/tgpost"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Daemon struct {
	Config *config.Config
	Ticker *time.Ticker
	Store  *store.Store
}

func New(configPath string) *Daemon {
	c := config.New(configPath)
	db, err := store.New(c)
	if err != nil {
		log.Fatalln("Ошибка при создании подключения к БД.")
	}

	return &Daemon{
		Config: c,
		Ticker: time.NewTicker(time.Duration(c.MainTickerTime) * time.Millisecond),
		Store:  db,
	}
}

func (d *Daemon) MotionHandler(e *rpidetectormongo.Event) (int, error) {
	err := tgpost.SendText(e.GetMotionMessage())
	if err != nil {
		return tgpost.StatusNotSent, errors.New("ошибка отправки текста о срабатывании")
	}
	err = d.Store.Collection.Update(bson.M{"_id": e.ID}, e)
	if err != nil {
		return store.StatusNotUpdated, errors.New(fmt.Sprintf("ошибка обновления записи в БД, id записи: %s", e.ID.Hex()))
	}

	return store.StatusUpdated, nil
}
