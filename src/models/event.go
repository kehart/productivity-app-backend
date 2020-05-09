package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Implements IEvent
type SleepEvent struct {
	Id 				primitive.ObjectID  `json:"id" bson:"_id"`
	UserId 			primitive.ObjectID  `json:"user_id" bson:"user_id" valid:"type(mongoid)"`
	Type 			string 				`json:"type" bson:"type"`
	SleepTime 		time.Time	`json:"sleep_time" bson:"sleep_time" valid:"rfc3339"` // maybe change to rfc3339WithoutZone
	WakeupTime 		time.Time	`json:"wakeup_time" bson:"wakeup_time" valid:"rfc3339"`
	/* Below are OPTIONAL fields */
	WakeupFeeling 	string		`json:"wakeup_feeling" bson:"wakeup_feeling" valid:"type(string), optional"` // custom; one of enum
	SleepFeeling 	string		`json:"sleep_feeling" bson:"sleep_feeling" valid:"type(string), optional"` // custom: one of enum
	QualityOfSleep 	int 		`json:"quality_of_sleep" bson:"quality_of_sleep" valid:"type(int), optional"` // [1, 10]
	AlarmUsed 		int		`json:"alarm_used" bson:"alarm_used" valid:"type(int), optional"`
	OwnBed 			int		`json:"own_bed" bson:"own_bed" valid:"type(int), optional"`
}

/*
Valid Quality of Sleep is in [1,10], -1 indicates no response
AlarmUsed and OwnBed use [0, 1] as true/false, -1 indicates no response
*/

/*
IEvent Implementation
*/

func (se SleepEvent) GetType() string {
	return se.Type
}

func (se SleepEvent) GetUserId() primitive.ObjectID {
	return se.UserId
}

func (se SleepEvent) GetId() primitive.ObjectID {
	return se.Id
}

// Semantic Validation of Certain Field Properties
func (se SleepEvent) Validate() error {
	if !( (1 <= se.QualityOfSleep && se.QualityOfSleep <= 10) || se.QualityOfSleep == -1) {
		err := errors.New("invalid quality_of_sleep value")
		return err
	}
	if !(-1 <= se.AlarmUsed && se.AlarmUsed <= 1) {
		err := errors.New("invalid alarm_used value")
		return err
	}
	if !(-1 <= se.OwnBed && se.OwnBed <= 1) {
		err := errors.New("invalid own_bed value")
		return err
	}
	return nil
}

type FoodData struct {
	FoodItem	string		`json:"food_item" bson:"food_item"`
	Quantity	string 		`json:"quantity" bson:"quantity"` // [item:float64][unit:enum]
}

// Implements IEvent
type DietEvent struct {
	Id 				primitive.ObjectID  `json:"id" bson:"_id"`
	UserId 			primitive.ObjectID  `json:"user_id" bson:"user_id" valid:"type(mongoid)"`
	Type 			string 				`json:"type" bson:"type"`
	// other custom fields
	TimeEaten		time.Time			`json:"time_eaten" bson:"time_eaten" valid:"rfc3339"`
	Items			[]FoodData			`json:"items" bson:"items"`
	Feeling			string				`json:"feeling" bson:"feeling" valid:"type(string), optional"` // custom; one of enum
}

/*
IEvent Implementation
 */

func (de DietEvent) GetType() string {
	return de.Type
}

func (de DietEvent) GetUserId() primitive.ObjectID {
	return de.UserId
}

func (de DietEvent) GetId() primitive.ObjectID {
	return de.Id
}

// Semantic Validation of Certain Field Properties
func (de DietEvent) Validate() error {
	// TODO
	return nil
}








