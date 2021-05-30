package servicecontainer_test

import (
	"github.com/vasilpatelnya/rpi-home/container/servicecontainer"
)

func getTestServiceContainer() servicecontainer.ServiceContainer {
	sc := servicecontainer.ServiceContainer{}
	err := sc.InitApp()
	if err != nil {
		sc.Logger.Fatalf("create test service container: fail. %s", err.Error())
	}

	return sc
}
