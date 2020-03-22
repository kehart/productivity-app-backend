package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/productivity-app-backend/src/managers"
	"github.com/productivity-app-backend/src/utils"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)



type SleepHandler struct {
	SleepManager managers.SleepManager
}

func (sh SleepHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
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

	err := sh.SleepManager.InsertSleepEvent(&newSleepEvent); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSleepEvent)
}

func (sh SleepHandler) GetAllEventsByType() {
	// get all events with type = sleep
	fmt.Println("LOG: GetAllEventsByType called")


}

func (sh SleepHandler) GetEventById() {

}

// get all sleep events for user
// allow parameters for time, ownBed, etc.

// Idea;
/* Have the generic event handler and pawn it off to the appropriate type manager

 */