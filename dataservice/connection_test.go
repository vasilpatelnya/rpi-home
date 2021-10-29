package dataservice_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vasilpatelnya/rpi-home/dataservice"
	"github.com/vasilpatelnya/rpi-home/tool/testhelpers"
	"testing"
)

func TestAssertCreateConnectionContainer(t *testing.T) {
	sc, err := testhelpers.GetTestContainer()
	assert.Nil(t, err)
	cc, err := dataservice.NewConnectionContainer(sc.AppConfig.Database)
	assert.Nil(t, err)
	t.Logf("%+v", cc)
}
