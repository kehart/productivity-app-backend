package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	FirstName 	string `json:"first_name" bson:"first_name" valid:"type(string)"`
	LastName  	string `json:"last_name" bson:"last_name" valid:"type(string)"`
	ID			primitive.ObjectID `json:"id" bson:"_id" valid:"-"`
}


type GoalCategory string
const (
	Sleep	GoalCategory = "sleep"
)

type GoalType string
const (
	HoursSlept	GoalType = "hours_slept"
)

// How to deal with enums https://www.ribice.ba/golang-enums/

type Goal struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" valid:"-"`
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id" valid:"required"` // valid:"type(mongoid)
	GoalCategory GoalCategory       `json:"goal_category" bson:"goal_category" valid:"required"` //valid:"type(string)"
	GoalType     GoalType           `json:"goal_name" bson:"goal_name" valid:"required"` // valid:"type(string)"
	TargetValue  interface{}        `json:"target_value" bson:"target_value" valid:"required"`
}


type Feeling string

const (
	Sad 		Feeling = "sad"
	Happy 		Feeling = "happy"
	Tired		Feeling = "tired"
	Anxious 	Feeling = "anxious"
	Refreshed 	Feeling = "refreshed"
	Excited 	Feeling = "excited"
)

type SleepEvent struct {
	Id 		primitive.ObjectID  `json:"id" bson:"_id"`
	UserId 	primitive.ObjectID  `json:"user_id" bson:"user_id" valid:"type(mongoid)"`
	SleepTime 		time.Time	`json:"sleep_time" bson:"sleep_time" valid:"rfc3339"` // maybe change to rfc3339WithoutZone
	WakeupTime 		time.Time	`json:"wakeup_time" bson:"wakeup_time" valid:"rfc3339"`
	WakeupFeeling 	Feeling		`json:"wakeup_feeling" bson:"wakeup_feeling" valid:"type(string), optional"` // custom; one of enum
	SleepFeeling 	Feeling		`json:"sleep_feeling" bson:"sleep_feeling" valid:"type(string), optional"` // custom: one of enum
	QualityOfSleep 	int 		`json:"quality_of_sleep" bson:"quality_of_sleep" valid:"type(itn), optional"` // [1, 10]
	AlarmUsed 		bool		`json:"alarm_used" bson:"alarm_used" valid:"type(bool), optional"`
	OwnBed 			bool		`json:"own_bed" bson:"own_bed" valid:"type(bool), optional"`
}


type (
	// db abstraction [users so far]
	Store interface {
		Create(user *User) error
		Delete(Id primitive.ObjectID) error
		FindById(id primitive.ObjectID) (*User, error) // done
		FindAll() (*[]User, error) // done
		Update(id primitive.ObjectID, user *User) (*User, error) // done
	}

	// The app 'context'; not sure if i need this
	//app struct {
	//	MongoDb store
	//}

	// Other types which are the models
)
