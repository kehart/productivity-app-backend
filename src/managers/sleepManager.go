package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"net/http"
)

type SleepManager struct {
	Session *mgo.Session
}

func (sm SleepManager) InsertSleepEvent(newSleepEvent *utils.SleepEvent) *utils.HTTPErrorLong {
	fmt.Println("LOG: Manager.InsertSleepEvent called")

	// Assign new ID to new user
	newSleepEvent.Id = primitive.NewObjectID()
	var user utils.User
	err := sm.Session.DB(utils.DbName).C(utils.UserCollection).FindId(newSleepEvent.UserId).One(&user); if err != nil {
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
	errLong := sm.Session.DB(utils.DbName).C(utils.EventCollection).Insert(newSleepEvent); if errLong != nil {
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