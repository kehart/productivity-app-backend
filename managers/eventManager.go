package managers

import (
	"fmt"
	"github.com/productivity-app-backend/interfaces"
	"github.com/productivity-app-backend/models"
	"github.com/productivity-app-backend/utils"
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
	userId := (*event).GetUserId()
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

	var results []map[string]interface{}
	var err interface{} // change
	if queryVals != nil {
		finalQueryVals := utils.ParseQueryString(queryVals)
		err = em.Store.FindAll(utils.EventCollection, &results, finalQueryVals)
	} else {
		err = em.Store.FindAll(utils.EventCollection, &results)
	}

	// TODO probably want to provide parallelism here
	var events []interfaces.IEvent
	fmt.Println(results)
	for _, e := range results {
		fmt.Println(e)
		event, err := interfaces.NewEventCreated(e); if err != nil {
			fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
			log.Println(utils.ErrorLog + err.Error()) // TODO ??
			return nil, &fullErr
		}
		events = append(events, event)
	}

	if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &events, nil
}

func (em EventManagerImpl) GetSingleEvent(objId primitive.ObjectID) (*interfaces.IEvent, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "EventManager:GetSingleEvent called")

	var event map[string]interface{}
	fmt.Println(objId)
	err := em.Store.FindById(objId, utils.EventCollection, &event); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), utils.NotFoundErrorString("Event", objId.String()), http.StatusNotFound)
		log.Println(utils.ErrorLog + err.Error()) // TODO ??
		return nil, &fullErr
	}
	specificEvent, err := interfaces.NewEventCreated(event); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + err.Error()) // TODO ??
		return nil, &fullErr
	}
	return &specificEvent, nil
}