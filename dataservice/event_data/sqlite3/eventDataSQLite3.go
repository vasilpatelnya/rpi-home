package sqlite3

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vasilpatelnya/rpi-home/model"
	"time"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type EventDataSQLite3 struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func (data *EventDataSQLite3) GetAllByStatus(s int) ([]model.Event, error) {
	var events []model.Event

	q := fmt.Sprintf("select * from 'events' where status = %d", s)
	rows, err := data.DB.Query(q)
	if err != nil {
		panic(err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		p := model.Event{}
		err = rows.Scan(&p.SqlID, &p.Status, &p.Type, &p.Device, &p.Name, &p.Updated, &p.Created)
		if err != nil {
			data.Logger.Errorf("Ошибка при получении событий по статусу '%d': %s", s, err.Error())

			return nil, err
		}
		events = append(events, p)
	}

	return events, nil
}

func (data *EventDataSQLite3) Save(e *model.Event) error {
	q := fmt.Sprintf(
		`INSERT INTO events (status, type, device, name, updated, created) 
		VALUES(%d, %d, '%s', '%s', %d, %d)`,
		e.Status, e.Type, e.Device, e.Name, e.Updated, e.Created,
	)
	if e.SqlID != 0 {
		q = fmt.Sprintf(
			`UPDATE events
			SET status = %d, type = %d, device = '%s', name = '%s', created = %d, updated = %d
			WHERE id = %d`, e.Status, e.Type, e.Device, e.Name, e.Created, e.Updated, e.SqlID)
	}

	_, err := data.DB.Exec(q)

	if err != nil {
		data.Logger.Errorf("ошибка сохранения неотправленного события: %s", err.Error())

		return err
	}

	return nil
}

func (data *EventDataSQLite3) SaveUpdated(e *model.Event, status int) error {
	if e.SqlID == 0 {
		return errors.New("event without id")
	}

	e.Status = status
	e.Updated = time.Now().UnixNano()

	return data.Save(e)
}
