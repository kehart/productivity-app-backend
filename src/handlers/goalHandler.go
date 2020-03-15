package handlers


import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	GoalCollection = "goals"
)

type GoalHandler struct {
	Session *mgo.Session
}

type goal struct {
	ID 				primitive.ObjectID `json:"id" bson:"_id"`
	UserId			primitive.ObjectID `json:"user_id" bson:"user_id"`
	GoalCategory	string `json:"goal_category" bson:"goal_category"` // always 'sleep' for now
	GoalName		string `json:"goal_name" bson:"goal_name"`
	TargetValue		int `json:"target_value" bson:"target_value"` // probably change this type later or make abstract or something
}

// Uses reqBody to create a new goal and inserts into DB
func (gh GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createGoal called")

	// Read request
	var newGoal goal
	reqBody, err := ioutil.ReadAll(r.Body); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	json.Unmarshal(reqBody, &newGoal)

	// Validate user_id passed in
	count, err := gh.Session.DB("admin-db").C(UserCollection).FindId(newGoal.UserId).Count(); if count != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("user not found")
		return
	}

	// Insert goal into db
	newGoal.ID = primitive.NewObjectID()
	err = gh.Session.DB("admin-db").C(GoalCollection).Insert(newGoal); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	objId, err := primitive.ObjectIDFromHex(goalID); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	var goal goal

	_, err = ioutil.ReadAll(r.Body); if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// search for user
	err = gh.Session.DB("admin-db").C(GoalCollection).FindId(objId).One(&goal); if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(goal)
	w.WriteHeader(http.StatusOK) // TODO: superflous?
}

// Returns list of all goals in DB // TODO change this to all ACTIVE goals
func (gh GoalHandler) GetGoals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getAllGoals called")

	var results []goal
	err := gh.Session.DB("admin-db").C(GoalCollection).Find(nil).All(&results); if err != nil {
		// TODO: what should actually happen here?
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: what do i return if empty
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}


func (gh GoalHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	// probably don't want to udpate; soft delete the old one and create a new 'updated record'
}


func (gh GoalHandler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	// probably a soft-delete
}