package models

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

/* Goal Type Definitions */
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


// TODO change to have this as a receiver?
func ValidateGoal(goal *Goal) *utils.HTTPErrorLong {
	if !goal.GoalCategory.isValid(goal.GoalType) {
		errBody := utils.HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, goal_category and goal_name should be a valid pair",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	} else if !goal.GoalType.isValid(goal.TargetValue) {
		errBody := utils.HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, the type of target_value does not match goal_name",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	}
	return nil
}


// Validate goal type with respect to goal category
func (gc GoalCategory) isValid(gt GoalType) bool {
	switch gc {
	case Sleep:
		if gt == HoursSlept {
			return true
		}
		return false
	}
	return false
}

// Validate target value with respect to goal type
func (gt GoalType) isValid(target interface{}) bool {
	switch gt {
	case HoursSlept:
		tType := fmt.Sprintf("%T", target)
		if tType =="int" || tType == "float64" {
			return true
		}
		fmt.Println(tType)
		return false
	}
	return false
}