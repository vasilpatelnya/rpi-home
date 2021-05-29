package sqlite3_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/config"
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
	"github.com/vasilpatelnya/rpi-home/model"
	"os"
	"testing"
)

var testContainer servicecontainer.ServiceContainer

func TestMain(m *testing.M) {
	testContainer = getTestServiceContainer()
	_, db := testContainer.DB.SQLite3.C()

	_, err := db.Exec(`DROP TABLE IF EXISTS 'events';`)
	if err != nil {
		testContainer.Logger.Fatalf("Drop table error: %s", err.Error())
	}

	_, err = db.Exec(`CREATE TABLE 'events' (
        'id' INTEGER PRIMARY KEY AUTOINCREMENT,
        'type' INTEGER NULL,
        'status' INTEGER NULL,
        'device' VARCHAR(64) NULL,
        'name' VARCHAR(64) NULL,
        'created' INTEGER NULL,
        'updated' INTEGER NULL);`)
	if err != nil {
		testContainer.Logger.Fatalf("Migration error: %s", err.Error())
	}
	os.Exit(m.Run())
}

func getTestServiceContainer() servicecontainer.ServiceContainer {
	sc := servicecontainer.ServiceContainer{}
	err := sc.InitApp()
	if err != nil {
		sc.Logger.Fatalf("create test service container: fail. %s", err.Error())
	}
	sc.Repo = servicecontainer.GetRepo(sc.DB.SQLite3, sc.Logger)

	return sc
}

func getAllEvents(t *testing.T, conn *config.SQLite3Connection) []model.Event {
	_, db := conn.C()
	rows, err := db.Query("SELECT * FROM events WHERE id = 1")
	assert.Nil(t, err)
	defer func() { _ = rows.Close() }()

	events := make([]model.Event, 0)
	for rows.Next() {
		p := model.Event{}
		err = rows.Scan(&p.SqlID, &p.Status, &p.Type, &p.Device, &p.Name, &p.Updated, &p.Created)
		assert.Nil(t, err)
		events = append(events, p)
	}

	return events
}

func clearTable(t *testing.T, conn *config.SQLite3Connection, table string) {
	_, db := conn.C()
	_, err := db.Exec(`DELETE FROM ` + table)
	assert.Nil(t, err)
}
