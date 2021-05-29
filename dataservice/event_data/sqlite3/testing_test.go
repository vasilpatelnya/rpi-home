package sqlite3_test

import (
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
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
        'name' VARCHAR(64) NULL,
        'device' VARCHAR(64) NULL,
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
