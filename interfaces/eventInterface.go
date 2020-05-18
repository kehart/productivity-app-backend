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
	log.Println(utils.InfoLog + "EventInterface::NewEvent called")

	eventType := json["type"]
	switch eventType{
	case "sleep":
		return NewSleepEvent(json)
	case "diet":
		return NewDietEvent(json)
	default:
		return nil, errors.New("error type not defined")
	}
}

func NewEventCreated(bson map[string]interface{}) (IEvent, error) {
	log.Print(utils.InfoLog + "EventInterface::NewEventCreated called")
	eventType := bson["type"]
	switch eventType{
	case "sleep":
		return NewSleepEventCreated(bson)
	case "diet:":
		return NewDietEventCreated(bson)
	default:
		return nil, errors.New("error type not defined")
	}
}


/*
Custom Constructor
*/

func NewSleepEvent(json map[string]interface{}) (*models.SleepEvent, error) {
	log.Println(utils.InfoLog + "EventInterface::NewSleepEvent called")

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
		se.Id.UnmarshalJSON(id.([]byte))
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

/*
Custom Constructor
*/

func NewDietEvent(json map[string]interface{}) (*models.DietEvent, error) {
	log.Println(utils.InfoLog + "EventInterface::NewDietEvent called")

	var de models.DietEvent

	// Mandatory fields
	uid := json["user_id"]; if uid != nil {
		objId, e := primitive.ObjectIDFromHex(uid.(string)); if e != nil {
			err := errors.New("error parsing user_id string")
			return nil, err
		}
		de.UserId = objId
	} else {
		err := errors.New("no user_id given")
		return nil, err
	}
	eType := json["type"]; if eType != nil {
		de.Type = eType.(string)
	} else {
		err := errors.New("no type given")
		return nil, err
	}

	st := json["time_eaten"]; if st != nil {
		t, e := time.Parse(time.RFC3339, st.(string)); if e != nil {
			err := errors.New("error parsing time_eaten string")
			return nil, err
		}
		de.TimeEaten = t
	} else {
		err := errors.New("no time_eaten given")
		return nil, err
	}
	de.Id = primitive.NewObjectID()

	items := json["items"]; if items != nil{
		foodItems, err := ParseFoodItems(items.([]interface{})); if err != nil {
			return nil, err
		} else {
			de.Items = foodItems
		}
	} else {
		err := errors.New("no items given")
		return nil, err
	}

	// Optional Fields
	wf := json["feeling"]; if wf != nil {
		de.Feeling = wf.(string)
	}
	return &de, nil
}


func NewDietEventCreated(bsonMap map[string]interface{}) (*models.DietEvent, error) {
	log.Println(utils.InfoLog + "EventInterface::NewDietEventCreated called")

	var de models.DietEvent

	uid := bsonMap["user_id"]; if uid != nil {
		de.UserId.UnmarshalJSON(uid.([]byte))
	} else {
		err := errors.New("no user_id given")
		return nil, err
	}
	eType := bsonMap["type"]; if eType != nil {
		de.Type = eType.(string)
	} else {
		err := errors.New("no type given")
		return nil, err
	}
	id := bsonMap["_id"]; if id != nil {
		de.Id.UnmarshalJSON(id.([]byte))
	} else {
		err := errors.New("no id given")
		return nil, err
	}

	it := bsonMap["items"]; if it != nil {
		de.Items = it.([]models.FoodData)
	} else {
		err := errors.New("no item given")
		return nil, err
	}

	wt := bsonMap["time_eaten"]; if wt != nil {
		de.TimeEaten = wt.(time.Time)
	} else {
		err := errors.New("no time_eaten given")
		return nil, err
	}

	// Optional Fields
	wf := bsonMap["feeling"]; if wf != nil {
		de.Feeling = wf.(string)
	}
	return &de, nil
}

func ParseFoodItems(items []interface{}) ([]models.FoodData, error) {
	finalList := []models.FoodData{}
	for _, e := range items {
		var item models.FoodData
		itemMap := e.(map[string]interface{})
		fi := itemMap["food_item"]; if fi != nil {
			item.FoodItem = fi.(string)
		} else {
			err := errors.New("no food_item given in items")
			return nil, err
		}
		qty := itemMap["quantity"]; if qty != nil {
			item.Quantity = qty.(string)
		} else {
			err := errors.New("no quantity given in items")
			return nil, err
		}
		finalList = append(finalList, item)
	}
	return finalList, nil
}