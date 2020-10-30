package dataservice

import (
	"github.com/vasilpatelnya/rpi-home/model"
)

type EventData interface {
	Save(e *model.Event) error
	SaveUpdated(e *model.Event, status int) error
	GetAllByStatus(s int) ([]model.Event, error)
}
