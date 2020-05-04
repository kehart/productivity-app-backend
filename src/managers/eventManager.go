package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
)

type EventManagerImpl struct {
	Store 	interfaces.Store
}

func (em EventManagerImpl) CreateEvent(event *interfaces.IEvent) (*interfaces.IEvent, *models.HTTPErrorLong){
	fmt.Println("LOG: Manager.InsertEvent called")

	// Validate user being referenced exists
	userId := (*event).GetUnderlyingEvent().UserId
	var user models.User
	err := em.Store.FindById2(userId, utils.UserCollection, &user); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id user_id not found", // TODO
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}

	// Insert user into DB
	errLong := em.Store.Create(event, utils.EventCollection); if errLong != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return event, nil
}

func (em EventManagerImpl) GetEvents(queryVals *url.Values) (*[]interfaces.IEvent, *models.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetEvents called")

	var results []interfaces.IEvent // change
	var err interface{} // change
	if queryVals != nil {
		finalQueryVals := utils.ParseQueryString(queryVals)
		err = em.Store.FindAll2(utils.EventCollection, &results, finalQueryVals)
	} else {
		err = em.Store.FindAll2(utils.EventCollection, &results)
	}

	if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return &results, nil
}

func (em EventManagerImpl) GetSingleEvent(objId primitive.ObjectID) (*interfaces.IEvent, *models.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetSingleEvent called")

	var event interfaces.IEvent
	err := em.Store.FindById2(objId, utils.EventCollection, &event); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("Event with id %s not found", objId.String()),
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return &event, nil
}