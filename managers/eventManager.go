package managers

import (
	"context"
	"fmt"
	"github.com/productivity-app-backend/interfaces"
	"github.com/productivity-app-backend/models"
	"github.com/productivity-app-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		log.Println(utils.ErrorLog + err.Error())
		return nil, &fullErr
	}

	// Insert user into DB
	eventMap := (*event).ToMap()
	errLong := em.Store.Create(eventMap, utils.EventCollection); if errLong != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + err.Error())
		return nil, &fullErr
	}
	return event, nil
}

func (em EventManagerImpl) GetEvents(queryVals *url.Values) (*[]interfaces.IEvent, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "EventManager:GetEvents called")

	var events []interfaces.IEvent
	decoder := func (cur *mongo.Cursor) error {
		for cur.Next(context.TODO()) {
			var eventMap map[string]interface{}
			err := cur.Decode(&eventMap); if err != nil {
				return  err
			}
			event, err := interfaces.NewEventCreated(eventMap); if err != nil {
				return err
			}
			events = append(events, event)
		}
		err := cur.Err()
		return err
	}

	var err interface{} // change
	if queryVals != nil {
		finalQueryVals := utils.ParseQueryString(queryVals)
		err = em.Store.FindAll(utils.EventCollection, nil, decoder, finalQueryVals)
	} else {
		err = em.Store.FindAll(utils.EventCollection, nil, decoder)
	}

	if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + fullErr.Error.ErrorMessage.(string))
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