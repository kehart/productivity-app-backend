package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/url"
)

type EventManager struct {
	Session *mgo.Session
}

func (em EventManager) InsertEvent(newSleepEvent *utils.SleepEvent) *utils.HTTPErrorLong {
	fmt.Println("LOG: Manager.InsertEvent called")

	// Assign new ID to new user
	newSleepEvent.Id = primitive.NewObjectID()
	var user utils.User
	err := em.Session.DB(utils.DbName).C(utils.UserCollection).FindId(newSleepEvent.UserId).One(&user); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id user_id not found", // TODO
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return &fullErr
	}
	// Insert user into DB
	errLong := em.Session.DB(utils.DbName).C(utils.EventCollection).Insert(newSleepEvent); if errLong != nil {
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

func (em EventManager) GetEvents(queryVals *url.Values) (*[]utils.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.GetGoals called")

	finalQueryVals := utils.ParseQueryString(queryVals)
	var results []utils.Goal
	err := em.Session.DB(utils.DbName).C(utils.EventCollection).Find(&finalQueryVals).All(&results); if err != nil {
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

// TODO this should not return a SleepEvent but a generic event
func (em EventManager) GetSingleEvent(objId primitive.ObjectID) (*utils.SleepEvent, *utils.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetSingleEvent called")

	var event utils.SleepEvent
	err := em.Session.DB(utils.DbName).C(utils.EventCollection).FindId(objId).One(&event); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"Event with id ID not found", // TODO figure out string interpolation
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return &event, nil
}