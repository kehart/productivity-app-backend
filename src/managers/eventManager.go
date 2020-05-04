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
	err := em.Store.FindById(userId, utils.UserCollection, &user); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), "User with id user_id not found", http.StatusNotFound)
		return nil, &fullErr
	}

	// Insert user into DB
	errLong := em.Store.Create(event, utils.EventCollection); if errLong != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
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
		err = em.Store.FindAll(utils.EventCollection, &results, finalQueryVals)
	} else {
		err = em.Store.FindAll(utils.EventCollection, &results)
	}

	if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		return nil, &fullErr
	}
	return &results, nil
}

func (em EventManagerImpl) GetSingleEvent(objId primitive.ObjectID) (*interfaces.IEvent, *models.HTTPErrorLong) {
	fmt.Println("LOG: EventManager.GetSingleEvent called")

	var event interfaces.IEvent
	err := em.Store.FindById(objId, utils.EventCollection, &event); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), fmt.Sprintf("Event with id %s not found", objId.String()), http.StatusNotFound)
		return nil, &fullErr
	}
	return &event, nil
}