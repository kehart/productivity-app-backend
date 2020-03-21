package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"net/http"
)

type GoalManager struct {
	Session *mgo.Session
}

func (gm GoalManager) CreateGoal(newGoal *utils.Goal) (*utils.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.CreateGoal called")

	// Check the userId in newGoal exists
	count, err := gm.Session.DB(utils.DbName).C(utils.UserCollection).FindId(newGoal.UserId).Count(); if count != 1 {
		errBody := utils.HttpError{
			ErrorCode: 		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	"user with id user_id not found",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return nil, &fullErr
	}

	// Insert goal into db
	newGoal.ID = primitive.NewObjectID()
	err = gm.Session.DB(utils.DbName).C(utils.GoalCollection).Insert(newGoal); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"internal server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return newGoal, nil
}

func (gm GoalManager) GetSingleGoal(objId primitive.ObjectID) (*utils.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.GetSingleGoal called")

	var goal utils.Goal
	err := gm.Session.DB(utils.DbName).C(utils.GoalCollection).FindId(objId).One(&goal); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"Goal with id ID not found", // TODO figure out string interpolation
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return &goal, nil
}

func (gm GoalManager) GetGoals() (*[]utils.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.GetGoals called")

	var results []utils.Goal
	err := gm.Session.DB(utils.DbName).C(utils.GoalCollection).Find(nil).All(&results); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return &results, nil
}

func (gm GoalManager) DeleteGoal(objId primitive.ObjectID) *utils.HTTPErrorLong {
	fmt.Println("LOG: GoalManager.DeleteGoal called")

	err := gm.Session.DB(utils.DbName).C(utils.GoalCollection).RemoveId(objId); if err != nil {
		if err.Error() == "not found" {
			errBody := utils.HttpError{
				ErrorCode:		http.StatusText(http.StatusNotFound),
				ErrorMessage: 	"ID not found",
			}
			fullErr := utils.HTTPErrorLong{
				Error:      errBody,
				StatusCode: http.StatusNotFound,
			}
			return &fullErr
		}
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return &fullErr
	}
	return nil
}