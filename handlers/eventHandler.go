package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/interfaces"
	"github.com/productivity-app-backend/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)


type EventHandler struct {
	EventManager interfaces.IEventManager
}


func (eh EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "EventHandler:CreateEvent called")

	var eventMap map[string]interface{}

	reqBody, genErr := ioutil.ReadAll(r.Body); if genErr != nil {
		utils.ReturnWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), genErr.Error())
		return
	}
	json.Unmarshal(reqBody, &eventMap)

	// Custom Unmarshalling to Specific Event Object
	event, err := interfaces.NewEvent(eventMap); if err != nil {
		utils.ReturnWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}

	// Validate Event Object
	err = event.Validate(); if err != nil {
		utils.ReturnWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}

	// Insert the object
	createdEvent, longErr := eh.EventManager.CreateEvent(&event); if longErr != nil {
		utils.ReturnWithErrorLong(w, *longErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

// Use cases: get all events for a certain type (where type will be a query value)
func (eh EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "EventHandler:GetEvents called")

	var queryVals *url.Values // base type *map[string][]string
	queryValMap := r.URL.Query(); if len(queryValMap) == 0 {
		queryVals = nil
	} else {
		queryVals = &queryValMap
	}

	results, err := eh.EventManager.GetEvents(queryVals); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

func (eh EventHandler) GetSingleEvent(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "EventHandler:GetSingleEvent called")

	eventID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(eventID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	event, err := eh.EventManager.GetSingleEvent(objId); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		return
	}
	json.NewEncoder(w).Encode(event)
	w.WriteHeader(http.StatusOK)
}


// Idea;
/* Have the generic event handler and pawn it off to the appropriate type manager

 */