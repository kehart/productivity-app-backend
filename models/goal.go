package models

import (
	"errors"
	"fmt"
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

func NewGoal(obj map[string]interface{}) (*Goal, error) {
	id := obj["_id"]; if id == nil {
		return nil, errors.New("no id")
	}
	uId := obj["user_id"]; if uId == nil {
		return nil, errors.New("no uId")
	}
	gc := obj["goal_category"]; if gc == nil {
		return nil, errors.New("no goalCateogry")
	}
	gt := obj["goal_name"]; if gt == nil {
		return nil, errors.New("no goalType")
	}

	fmt.Println("this is gt", gt)
	tv := obj["target_value"]; if tv == nil {
		return nil, errors.New("no targetValue")
	}
	goal := Goal{
		ID:           id.(primitive.ObjectID),
		UserId:       uId.(primitive.ObjectID),
		GoalCategory: gc.(GoalCategory),
		GoalType:     gt.(GoalType),
		TargetValue:  tv,
	}
	return &goal, nil
}


func (g *Goal) Validate() *HTTPErrorLong {
	if !g.GoalCategory.isValid(g.GoalType) {
		fullErr := NewHTTPErrorLong(http.StatusText(http.StatusBadRequest), "Error, goal_category and goal_name should be a valid pair", http.StatusBadRequest)
		return &fullErr
	} else if !g.GoalType.isValid(g.TargetValue) {
		fullErr := NewHTTPErrorLong(http.StatusText(http.StatusBadRequest), "Error, the type of target_value does not match goal_name", http.StatusBadRequest)
		return &fullErr
	}
	return nil
}

// Validate goal type with respect to goal category
func (gc *GoalCategory) isValid(gt GoalType) bool {
	switch *gc {
	case Sleep:
		if gt == HoursSlept {
			return true
		}
		return false
	}
	return false
}

// Validate target value with respect to goal type
func (gt *GoalType) isValid(target interface{}) bool {
	switch *gt {
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