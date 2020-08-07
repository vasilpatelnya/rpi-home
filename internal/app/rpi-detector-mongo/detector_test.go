package rpi_detector_mongo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/internal/app/config"
	"github.com/vasilpatelnya/rpi-home/internal/app/store"
	"github.com/vasilpatelnya/rpi-home/internal/app/tgpost"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	testDirPath = "../../../testDir"
	backupPath  = "../../../backup"
	configPath  = "./../../../configs/test.env"
)

type Stats struct {
	All        int
	RightFiles int
	WrongFiles int
}

func TestNew(t *testing.T) {
	e := New()
	assert.True(t, reflect.DeepEqual(e, &Event{
		Type:    TypeUndefined,
		Status:  StatusNew,
		Name:    "",
		Device:  "",
		Created: e.Created,
		Updated: e.Updated,
	}))
}

func TestEvent_GetAllByStatus(t *testing.T) {
	c := config.New(configPath)
	s, err := store.New(c)
	assert.Nil(t, err)
	_, err = GetAllByStatus(s.Collection, StatusNew)
	assert.Nil(t, err)
}

func TestEvent_Save(t *testing.T) {
	c := config.New(configPath)
	s, err := store.New(c)
	assert.Nil(t, err)
	event := &Event{
		ID:      bson.NewObjectId(),
		Type:    1,
		Status:  1,
		Name:    "test",
		Device:  "test",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	err = event.Save(s.Collection)
	assert.Nil(t, err)
}

func TestEvent_HandlerMotionReady(t *testing.T) {
	err := os.Setenv("FILE_EXTENSION", ".mp4")
	assert.Nil(t, err)
	err = os.Setenv("APP_MODE", config.AppTest)
	assert.Nil(t, err)
	event := &Event{
		ID:      bson.NewObjectId(),
		Type:    TypeMovieReady,
		Status:  StatusNew,
		Name:    "test",
		Device:  "test",
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}
	statsTestDirStart, err := getStats(testDirPath)
	assert.Nil(t, err)
	statsBackupDirStart, err := getStats(backupPath)
	assert.Nil(t, err)
	status, err := event.HandlerMotionReady(testDirPath, backupPath)
	assert.Nil(t, err)
	assert.Equal(t, status, tgpost.StatusSent)
	statsTestDirEnd, err := getStats(testDirPath)
	assert.Nil(t, err)
	statsBackupDirEnd, err := getStats(testDirPath)
	assert.Nil(t, err)
	assert.Equal(t, statsTestDirEnd.All, statsTestDirStart.All-statsTestDirStart.RightFiles)
	// todo здесь спотыкается тест - разобраться.
	assert.Equal(t, statsBackupDirEnd.All, statsBackupDirStart.All+statsTestDirStart.RightFiles)

	files, err := ioutil.ReadDir(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == os.Getenv("FILE_EXTENSION") {
			todayDir := tgpost.GetTodayDir()
			box, err := ioutil.ReadFile(backupPath + "/" + f.Name())
			if err != nil {
				log.Println("Ошибка при попытке прочитать файл:", f.Name(), err)
			}
			fp := fmt.Sprintf("%s/%s/%s", testDirPath, todayDir, f.Name())
			err = ioutil.WriteFile(fp, box, 0777)
			if err != nil {
				log.Println("Ошибка при попытке скопировать файл:", f.Name(), err)
			}
			err = os.Remove(backupPath + "/" + f.Name())
			if err != nil {
				log.Println("Ошибка при попытке удалить файл:", f.Name(), err)
			}
		}
	}
}

func getStats(path string) (*Stats, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	stats := &Stats{
		All:        len(files),
		RightFiles: 0,
		WrongFiles: 0,
	}
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == os.Getenv("FILE_EXTENSION") && f.Size() < MaxSize {
			stats.RightFiles++
			continue
		}
		stats.WrongFiles++
	}

	return stats, nil
}
