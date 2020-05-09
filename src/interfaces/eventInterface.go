package interfaces

import (
	"errors"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"time"
)

/*
IEvent is an interface that defines the behaviour
of each concrete Event type
*/


type IEvent interface {
	Validate() error
	GetType() string
	GetId() primitive.ObjectID
	GetUserId() primitive.ObjectID
}

// Factory method for creating new IEvents
func NewEvent(json map[string]interface{}) (IEvent, error) {
	eventType := json["type"]
	switch eventType{
	case "sleep":
		return NewSleepEvent(json)
	case "diet:":
		return nil, nil // TODO
	default:
		return nil, errors.New("error type not defined")
	}
}

func NewEventCreated(json map[string]interface{}) (IEvent, error) {
	log.Print(utils.InfoLog + "EventInterface::NewEventCreated called")
	eventType := json["type"]
	switch eventType{
	case "sleep":
		return NewSleepEventCreated(json)
	case "diet:":
		return nil, nil // TODO
	default:
		return nil, errors.New("error type not defined")
	}
}


/*
Custom Constructor
*/


func NewSleepEvent(json map[string]interface{}) (*models.SleepEvent, error) {
	var se models.SleepEvent

	// Mandatory fields
	uid := json["user_id"]; if uid != nil {
		objId, e := primitive.ObjectIDFromHex(uid.(string)); if e != nil {
			err := errors.New("error parsing user_id string")
			return nil, err
		}
		se.UserId = objId
	} else {
		err := errors.New("no user_id given")
		return nil, err
	}
	eType := json["type"]; if eType != nil {
		se.Type = eType.(string)
	} else {
		err := errors.New("no type given")
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
	se.Id = primitive.NewObjectID()

	// Optional Fields
	wf := json["wakeup_feeling"]; if wf != nil {
		se.WakeupFeeling = wf.(string)
	}
	sf := json["sleep_feeling"]; if sf != nil {
		se.SleepFeeling = sf.(string)
	}
	qos := json["quality_of_sleep"]; if qos != nil {
		qosInt, e := strconv.Atoi(qos.(string)); if e != nil {
			err := errors.New("error parsing quality_of_sleep")
			return nil, err
		}
		se.QualityOfSleep = qosInt
	} else {
		se.QualityOfSleep = -1
	}
	au := json["alarm_used"]; if au != nil {
		auInt, e := strconv.Atoi(au.(string)); if e != nil {
			err := errors.New("error parsing alarm_used")
			return nil, err
		}
		se.AlarmUsed = auInt
	} else {
		se.AlarmUsed = -1
	}
	ob := json["own_bed"]; if ob != nil {
		obInt, e := strconv.Atoi(ob.(string)); if e != nil {
			err := errors.New("error parsing own_bed")
			return nil, err
		}
		se.OwnBed = obInt
	} else {
		se.OwnBed = -1
	}
	return &se, nil
}


func NewSleepEventCreated(bsonMap map[string]interface{}) (*models.SleepEvent, error) {
	log.Println(utils.InfoLog + "EventInterface::NewSleepEventCreated called")

	var se models.SleepEvent

	uid := bsonMap["user_id"]; if uid != nil {
		se.UserId.UnmarshalJSON(uid.([]byte))
	} else {
		err := errors.New("no user_id given")
		return nil, err
	}
	eType := bsonMap["type"]; if eType != nil {
		se.Type = eType.(string)
	} else {
		err := errors.New("no type given")
		return nil, err
	}
	id := bsonMap["_id"]; if id != nil {
		se.Id.UnmarshalJSON(id.([]byte))//
	} else {
		err := errors.New("no id given")
		return nil, err
	}

	st := bsonMap["sleep_time"]; if st != nil {
		se.SleepTime = st.(time.Time)
	} else {
		err := errors.New("no sleep_time given")
		return nil, err
	}

	wt := bsonMap["wakeup_time"]; if wt != nil {
		se.WakeupTime = wt.(time.Time)
	} else {
		err := errors.New("no wakeup_time given")
		return nil, err
	}
	se.Id = primitive.NewObjectID()

	// Optional Fields
	wf := bsonMap["wakeup_feeling"]; if wf != nil {
		se.WakeupFeeling = wf.(string)
	}
	sf := bsonMap["sleep_feeling"]; if sf != nil {
		se.SleepFeeling = sf.(string)
	}
	qos := bsonMap["quality_of_sleep"]; if qos != nil {
		se.QualityOfSleep = qos.(int)
	}
	au := bsonMap["alarm_used"]; if au != nil {
		se.AlarmUsed = au.(int)
	}
	ob := bsonMap["own_bed"]; if ob != nil {
		se.OwnBed = ob.(int)
	}
	return &se, nil
}


