package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

/*
An Event is a type which consists of :
- a BaseEvent (generic information)
- custom fields specific to the event type (e.g. SleepEvent)
 */


/* Event Type Definitions */

type Feeling string

const (
	Sad 		Feeling = "sad"
	Happy 		Feeling = "happy"
	Tired		Feeling = "tired"
	Anxious 	Feeling = "anxious"
	Refreshed 	Feeling = "refreshed"
	Excited 	Feeling = "excited"
)


func (f Feeling) isValid() bool {
	switch f {
	case Sad, Happy, Tired, Anxious, Refreshed, Excited:
		return true
	}
	return false
}


// Component of ever concrete event type
type BaseEvent struct {
	Id 		primitive.ObjectID  `json:"id" bson:"_id"`
	UserId 	primitive.ObjectID  `json:"user_id" bson:"user_id" valid:"type(mongoid)"`
}

// Implements IEvent
type SleepEvent struct {
	BaseEvent 		BaseEvent
	SleepTime 		time.Time	`json:"sleep_time" bson:"sleep_time" valid:"rfc3339"` // maybe change to rfc3339WithoutZone
	WakeupTime 		time.Time	`json:"wakeup_time" bson:"wakeup_time" valid:"rfc3339"`
	WakeupFeeling 	Feeling		`json:"wakeup_feeling" bson:"wakeup_feeling" valid:"type(string), optional"` // custom; one of enum
	SleepFeeling 	Feeling		`json:"sleep_feeling" bson:"sleep_feeling" valid:"type(string), optional"` // custom: one of enum
	QualityOfSleep 	int 		`json:"quality_of_sleep" bson:"quality_of_sleep" valid:"type(itn), optional"` // [1, 10]
	AlarmUsed 		bool		`json:"alarm_used" bson:"alarm_used" valid:"type(bool), optional"`
	OwnBed 			bool		`json:"own_bed" bson:"own_bed" valid:"type(bool), optional"`
}

/*
IEvent Implementation
*/

func (se SleepEvent) GetUnderlyingEvent() BaseEvent {
	return se.BaseEvent
}

// SleepEvent implements IEvent
func (se SleepEvent) Validate() error {
	return nil
}

/*
Custom Constructor
*/

// A failure here means a 400
func NewSleepEvent(json map[string]interface{}) (*SleepEvent, error) {
	var se SleepEvent

	// Mandatory fields
	uid := json["user_id"]; if uid != nil {
		objId, e := primitive.ObjectIDFromHex(uid.(string)); if e != nil {
			err := errors.New("error parsing user_id string")
			return nil, err
		}
		se.BaseEvent.UserId = objId
	} else {
		err := errors.New("no user_id given")
		return nil, err
	}

	st := json["sleep_time"]; if st != nil {
		t, e := time.Parse(time.RFC3339, st.(string)); if e != nil {
			err := errors.New("error parsing sleep_time string")
			return nil, err
		}
		se.SleepTime = t
	} else {
		err := errors.New("no sleep_time given")
		return nil, err
	}

	wt := json["wakeup_time"]; if wt != nil {
		t, e := time.Parse(time.RFC3339, wt.(string)); if e != nil {
			err := errors.New("error parsing wakeup_time string")
			return nil, err
		}
		se.WakeupTime = t
	} else {
		err := errors.New("no wakeup_time given")
		return nil, err
	}
	se.BaseEvent.Id = primitive.NewObjectID()

	// Optional Fields
	// TODO fix serialization when values not present
	wf := json["wakeup_feeling"]; if wf != nil {
		se.WakeupFeeling = wf.(Feeling)
	}
	sf := json["sleep_feeling"]; if sf != nil {
		se.SleepFeeling = sf.(Feeling)
	}
	qos := json["quality_of_sleep"]; if qos != nil {
		qosInt, e := strconv.Atoi(qos.(string)); if e != nil {
			err := errors.New("error parsing quality_of_sleep")
			return nil, err
		}
		se.QualityOfSleep = qosInt
	}
	au := json["alarm_used"]; if au != nil {
		auBool, e := strconv.ParseBool(au.(string)); if e != nil {
			err := errors.New("error parsing alarm_used")
			return nil, err
		}
		se.AlarmUsed = auBool
	}
	ob := json["own_bed"]; if ob != nil {
		obBool, e := strconv.ParseBool(ob.(string)); if e != nil {
			err := errors.New("error parsing own_bed")
			return nil, err
		}
		se.OwnBed = obBool
	}
	return &se, nil
}


// Implements IEvent
type DietEvent struct {
	BaseEvent		BaseEvent
	// other custom fields
}

/*
IEvent Implementation TODO
 */









