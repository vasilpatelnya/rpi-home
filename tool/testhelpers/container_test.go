package testhelpers_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"os"
	"testing"
)

func TestGetTestDataDir(t *testing.T) {
	want := os.Getenv("GOPATH") + "/src/github.com/vasilpatelnya/rpi-home/testdata"
	got := testhelpers.GetTestDataDir()
	assert.Equal(t, want, got)
}
