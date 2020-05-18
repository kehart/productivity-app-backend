package handlers

import (
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/interfaces"
	"github.com/productivity-app-backend/models"
	"github.com/productivity-app-backend/utils"
	"io/ioutil"
	"log"
	"net/http"
)

type GoalHandler struct {
	GoalManager interfaces.IGoalManager
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


// Handles POST /goals
func (gh GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "GoalHandler:CreateGoal called")

	// Read request
	var newGoal models.Goal

	reqBody, genErr := ioutil.ReadAll(r.Body); if genErr != nil {
		utils.ReturnWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), genErr.Error())
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	json.Unmarshal(reqBody, &newGoal)

	// Validate Syntax
	_, genErr = valid.ValidateStruct(&newGoal) ; if genErr != nil {
		utils.ReturnWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), genErr.Error())
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	// Validate Semantics
	err := newGoal.Validate(); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	_, err = gh.GoalManager.CreateGoal(&newGoal); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGoal)
}

// Uses req body and id from path to read a single goal
func (gh GoalHandler) GetSingleGoal(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "GoalHandler:GetSingleGoal called")

	goalID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(goalID);  if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	goal, err := gh.GoalManager.GetSingleGoal(objId); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	json.NewEncoder(w).Encode(goal)
	w.WriteHeader(http.StatusOK)
}

// Returns list of all goals in DB // TODO change this to all ACTIVE goals
func (gh GoalHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "GoalHandler:GetGoals called")

	// Parse query string
	queryVals := r.URL.Query() // returns map[string][]string

	results, err := gh.GoalManager.GetGoals(&queryVals); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
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
	log.Print(utils.InfoLog + "GoalHandler:DeleteGoal called")

	goalID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(goalID);  if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	err = gh.GoalManager.DeleteGoal(objId); if err != nil {
		utils.ReturnWithErrorLong(w, *err)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
