package rpi_detector_mongo

import (
	"fmt"
	"github.com/vasilpatelnya/rpi-home/internal/app/tgpost"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		Type:    TypeUndefined,
		Status:  StatusNew,
		Name:    "",
		Device:  "",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
}

func GetAllByStatus(c *mgo.Collection, s int) ([]Event, error) {
	var events []Event
	if err := c.Find(bson.M{"status": s}).All(&events); err != nil {
		log.Fatal(err)

		return nil, err
	}

	return events, nil
}

func (e *Event) Save(c *mgo.Collection) error {
	_, err := c.Upsert(bson.M{"_id": e.ID}, e)
	if err != nil {
		log.Println("ошибка сохранения неотправленного события", err)

		return err
	}

	return nil
}

func (e *Event) GetMotionMessage() string {
	text := fmt.Sprintf("Устройство %s зафиксировало движение! ", e.Device)
	var add string
	if e.Status == -1 {
		add = "Это сообщение отправлено с задержкой. Точное время срабатывания: " + time.Unix(e.Updated, 0).Format(time.Stamp)
	}

	return text + add
}

func (e *Event) HandlerMotionReady(dirname string) (int, error) {
	l, err := tgpost.GetTodayFileList(dirname)
	if err != nil {
		log.Println("Ошибка получения списка файлов в директории:", err.Error())

		return tgpost.StatusNotSent, err
	}
	for _, f := range l {
		todayDir := tgpost.GetTodayDir()
		fp := fmt.Sprintf("%s/%s/%s", dirname, todayDir, f.Name())
		ext := filepath.Ext(f.Name())
		if ext == os.Getenv("FILE_EXTENSION") && f.Size() > 0 {
			err := tgpost.SendFile(fp, "Видео от "+time.Unix(e.Created/1000000000, 0).String())
			if err != nil {
				log.Println("Ошибка при попытке отправить видео", f.Name(), err)

				return tgpost.StatusNotSent, err
			}
			log.Printf("файл %s был отправлен в телеграм", fp)
			box, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Println("Ошибка при попытке прочитать файл:", f.Name(), err)

				return tgpost.StatusNotSent, err
			}
			err = ioutil.WriteFile("./backup/"+f.Name(), box, 0777)
			if err != nil {
				log.Println("Ошибка при попытке скопировать файл:", f.Name(), err)

				return tgpost.StatusNotSent, err
			}
			err = os.Remove(fp)
			if err != nil {
				log.Println("Ошибка при попытке удалить файл:", f.Name(), err)

				return tgpost.StatusNotSent, err
			}
		}
	}

	return tgpost.StatusSent, nil
}

func (e *Event) SaveUpdated(c *mgo.Collection, status int) error {
	e.Status = status
	e.Updated = time.Now().UnixNano()

	return e.Save(c)
}
