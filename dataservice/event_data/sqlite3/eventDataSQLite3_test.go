package sqlite3_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventDataSQLite3_Save(t *testing.T) {
	sc := getTestServiceContainer()
	connection, db := sc.DB.SQLite3.C()
	time.Sleep(1 * time.Second)
	assert.Nil(t, db.Ping())

	defer func() { _ = connection.Close() }()
}
