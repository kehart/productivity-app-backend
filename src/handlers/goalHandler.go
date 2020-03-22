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

type GoalHandler struct {
	GoalManager managers.GoalManager
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
func (gh GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createGoal called")

	// Read request
	var newGoal utils.Goal

	// Validate and unmarshal to newUser
	rules := govalidator.MapData{
		"user_id": []string{"required"},
		"goal_category": []string{"required"},
		"goal_name": []string{"required"},
		"target_value": []string{"required"},
	}
	opts := govalidator.Options{
		Data:            &newGoal,
		Request:         r,
		RequiredDefault: true, // idk what this does
		Rules:           rules,
	}
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

	_, err := gh.GoalManager.CreateGoal(&newGoal); if err != nil {
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

	results, err := gh.GoalManager.GetGoals(); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}
/* TODO
func (gh GoalHandler) GetGoalsForUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getGoalsForUser called")

	var results []goal
	queryStrVals := r.URL.Query() // probably should validate this
	userId := queryStrVals["user_id"][0]
	if  userId != "" {
		objId, err := primitive.ObjectIDFromHex(userId); if err != nil {
			// err
		}
		queryStrVals["user_id"] = []primitive.ObjectID {objId}
	}

	err := gh.Session.DB("admin-db").C(GoalCollection).Find(queryStrVals).All(&results); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}
*/

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
