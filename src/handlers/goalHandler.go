package handlers


import (
	"encoding/json"
	"fmt"
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


func (gh GoalHandler) GetSingleGoal(w http.ResponseWriter, r *http.Request) {

}


func (gh GoalHandler) getGoals(w http.ResponseWriter, r *http.Request) {

}


func (gh GoalHandler) updateGoal(w http.ResponseWriter, r *http.Request) {
	// probably don't want to udpate; soft delete the old one and create a new 'updated record'
}


func (gh GoalHandler) deleteGoal(w http.ResponseWriter, r *http.Request) {
	// probably a soft-delete
}