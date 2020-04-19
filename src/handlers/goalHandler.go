package handlers

import (
	"encoding/json"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/utils"
	"io/ioutil"
	"net/http"
)

type GoalHandler struct {
	GoalManager utils.IGoalManager
}

// TODO should extend above to have some sort of status? like deleted/active, achieved/in-progress
// TODO maybe time horizon? (3 weeks, ongoing, etc.)
// TODO when altering, you should only ever change the TargetValue

// Uses reqBody to create a new goal and inserts into DB
/* Cases:
-happy path :)
-bad id :)
-invalid fields (empty) :)
-invalid fields (dont match type of GoalCategory or GoalType) :( TODO
 */

/*	var newUser utils.User

	reqBody, genErr := ioutil.ReadAll(r.Body); if genErr != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	"Bad request",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	json.Unmarshal(reqBody, &newUser)
	_, genErr = valid.ValidateStruct(&newUser) ; if genErr != nil {
			errBody := utils.HttpError{
				ErrorCode:		http.StatusText(http.StatusBadRequest),
				ErrorMessage:	genErr,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errBody)
			return
	}
	err := utils.ValidateUser(&newUser); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}*/

// Handles POST /goals
func (gh GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createGoal called")

	// Read request
	var newGoal utils.Goal

	reqBody, genErr := ioutil.ReadAll(r.Body); if genErr != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	"Bad request",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	json.Unmarshal(reqBody, &newGoal)
	_, genErr = valid.ValidateStruct(&newGoal) ; if genErr != nil {
		fmt.Println(genErr)
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	genErr,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}
	err := utils.ValidateGoal(&newGoal); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	_, err = gh.GoalManager.CreateGoal(&newGoal); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGoal)
}

// Uses req body and id from path to read a single goal
func (gh GoalHandler) GetSingleGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getSingleGoal called")

	goalID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(goalID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	goal, err := gh.GoalManager.GetSingleGoal(objId); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(goal)
	w.WriteHeader(http.StatusOK)
}

// Returns list of all goals in DB // TODO change this to all ACTIVE goals
func (gh GoalHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getGoals called")

	// Parse query string
	queryVals := r.URL.Query() // returns map[string][]string

	results, err := gh.GoalManager.GetGoals(&queryVals); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

// For now there is no update permitted
// In the future, this should just report a status update, like active, completed, in progress, etc.
// Target value should not change, you should have to complete one goal and create another

// Hard-delete (all else should be update)
func (gh GoalHandler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: deleteGoal called")

	goalID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(goalID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	err = gh.GoalManager.DeleteGoal(objId); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
