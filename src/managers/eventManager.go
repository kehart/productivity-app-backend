package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
)

type EventManagerImpl struct {
	Store 	interfaces.Store
}

func (em EventManagerImpl) CreateEvent(event *interfaces.IEvent) (*interfaces.IEvent, *utils.HTTPErrorLong){
	fmt.Println("LOG: Manager.InsertEvent called")

	// Validate user being referenced exists
	userId := (*event).GetUnderlyingEvent().UserId
	_, err := em.Store.FindById(userId, utils.UserCollection); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id user_id not found", // TODO
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}

	// Insert user into DB
	errLong := em.Store.Create(event, utils.EventCollection); if errLong != nil {
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
	return event, nil
}

func (em EventManagerImpl) GetEvents(queryVals *url.Values) (*[]interfaces.IEvent, *utils.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetEvents called")

	var results interface{} // change
	var err interface{} // change
	if queryVals != nil {
		finalQueryVals := utils.ParseQueryString(queryVals)
		results, err = em.Store.FindAll(utils.EventCollection, finalQueryVals)
	} else {
		results, err = em.Store.FindAll(utils.EventCollection)
	}

	if err != nil {
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
	return results.(*[]interfaces.IEvent), nil
}

func (em EventManagerImpl) GetSingleEvent(objId primitive.ObjectID) (*interfaces.IEvent, *utils.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetSingleEvent called")

	event, err := em.Store.FindById(objId, utils.EventCollection); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("Event with id %s not found", objId.String()),
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return event.(*interfaces.IEvent), nil
}