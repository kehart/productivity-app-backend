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

type Goal struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id"`
	GoalCategory GoalCategory       `json:"goal_category" bson:"goal_category"`
	GoalType     GoalType           `json:"goal_name" bson:"goal_name"`
	TargetValue  interface{}        `json:"target_value" bson:"target_value"`
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
// TODO delete?
type BaseEvent struct {
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id"`
	// dates
	// goal category
}

type SleepEvent struct {
	Id 		primitive.ObjectID  `json:"id" bson:"_id"`
	UserId 	primitive.ObjectID  `json:"user_id" bson:"user_id"`
	SleepTime 		time.Time	`json:"sleep_time" bson:"sleep_time"`
	WakeupTime 		time.Time	`json:"wakeup_time" bson:"wakeup_time"`
	WakeupFeeling 	Feeling		`json:"wakeup_feeling" bson:"wakeup_feeling"`
	SleepFeeling 	Feeling		`json:"sleep_feeling" bson:"sleep_feeling"`
	QualityOfSleep 	int 		`json:"quality_of_sleep" bson:"quality_of_sleep"` // [1, 10]
	AlarmUsed 		bool		`json:"alarm_used" bson:"alarm_used"`
	OwnBed 			bool		`json:"own_bed" bson:"own_bed"`
}
