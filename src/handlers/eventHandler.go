package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/managers"
	"github.com/productivity-app-backend/src/utils"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)



type EventHandler struct {
	SleepManager managers.EventManager
}

func (sh EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createUser called")

	var newSleepEvent utils.SleepEvent

	// Validate and unmarshal to newSleepEvent
	rules := govalidator.MapData{
		"user_id": []string{"required"},
		"sleep_time": []string{"required"},
		"wakeup_time": []string{"required"},
		"wakeup_feeling": []string{},
		"sleep_feeling": []string{},
		"quality_of_sleep": []string{},
		"alarm_used": []string{},
		"own_bed": []string{},
	}
	opts := govalidator.Options{
		Data:            &newSleepEvent,
		Request:         r,
		RequiredDefault: true, // idk what this does
		Rules:           rules,
	}

	// got here
	v := govalidator.New(opts)
	e := v.ValidateJSON(); if len(e) > 0 {
		validationError := map[string]interface{}{"validationError": e}
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	validationError,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	err := sh.SleepManager.InsertEvent(&newSleepEvent); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSleepEvent)
}

// Use cases: get all events for a certain type (where type will be a query value)
func (sh EventHandler) GetAllEventsByType(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: GetAllEventsByType called")

	// Parse query string
	queryVals := r.URL.Query() // returns map[string][]string

	results, err := sh.SleepManager.GetEvents(&queryVals); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

func (sh EventHandler) GetSingleEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: GetEventById called")

	eventID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(eventID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	event, err := sh.SleepManager.GetSingleEvent(objId); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(event)
	w.WriteHeader(http.StatusOK)
}

// get all sleep events for user
// allow parameters for time, ownBed, etc.

// Idea;
/* Have the generic event handler and pawn it off to the appropriate type manager

 */