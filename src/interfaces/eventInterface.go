package interfaces

import (
	"errors"
	"github.com/productivity-app-backend/src/models"
)

/*
IEvent is an interface that defines the behaviour
of each concrete Event type
*/


type IEvent interface {
	Validate() error // TODO so far in validators
	GetUnderlyingEvent() models.BaseEvent
}

// Factory method for creating new IEvents
func NewEvent(json map[string]interface{}) (IEvent, error) {
	eventType := json["type"]
	switch eventType{
	case "sleep":
		return models.NewSleepEvent(json)
	case "diet:":
		return nil, nil // TODO
	default:
		return nil, errors.New("error type not defined")
	}
}

