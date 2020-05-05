package managers

import (
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/url"
)

type EventManagerImpl struct {
	Store 	interfaces.Store
}

func (em EventManagerImpl) CreateEvent(event *interfaces.IEvent) (*interfaces.IEvent, *models.HTTPErrorLong){
	log.Print(utils.InfoLog + "EventManager:CreateEvent called")

	// Validate that user being referenced exists
	userId := (*event).GetUnderlyingEvent().UserId
	var user models.User
	err := em.Store.FindById(userId, utils.UserCollection, &user); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), utils.NotFoundErrorString("User", userId.String()), http.StatusNotFound)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}

	// Insert user into DB
	errLong := em.Store.Create(event, utils.EventCollection); if errLong != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return event, nil
}

func (em EventManagerImpl) GetEvents(queryVals *url.Values) (*[]interfaces.IEvent, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "EventManager:GetEvents called")

	var results []interfaces.IEvent
	var err interface{} // change
	if queryVals != nil {
		finalQueryVals := utils.ParseQueryString(queryVals)
		err = em.Store.FindAll(utils.EventCollection, &results, finalQueryVals)
	} else {
		err = em.Store.FindAll(utils.EventCollection, &results)
	}

	if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &results, nil
}

func (em EventManagerImpl) GetSingleEvent(objId primitive.ObjectID) (*interfaces.IEvent, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "EventManager:GetSingleEvent called")

	var event interfaces.IEvent
	err := em.Store.FindById(objId, utils.EventCollection, &event); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), utils.NotFoundErrorString("Event", objId.String()), http.StatusNotFound)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &event, nil
}