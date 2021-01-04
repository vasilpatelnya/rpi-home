package fs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// GetTodayFileList ...
func GetTodayFileList(dirname, layout string) ([]os.FileInfo, error) {
	if Exists(dirname) {
		files, err := ioutil.ReadDir(GetTodayPath(dirname, layout))
		if err != nil {
			return nil, err
		}

		return files, nil
	}

	return nil, errors.New(dirname + " directory is not exist")
}

// GetTodayDir ...
func GetTodayDir(layout string) string {
	t := time.Now()

	return t.Format(layout)
}

// GetTodayPath ...
func GetTodayPath(dirname, layout string) string {
	return fmt.Sprintf("%s/%s", dirname, GetTodayDir(layout))
}
