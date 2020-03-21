package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	FirstName 	string `json:"first_name" bson:"first_name"`
	LastName  	string `json:"last_name" bson:"last_name"`
	ID			primitive.ObjectID `json:"id" bson:"_id"`
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
