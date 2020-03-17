package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/utils"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
)
const (
	UserCollection = "users"
)

type UserHandler struct {
	Session *mgo.Session
}

type user struct {
	FirstName 	string `json:"first_name" bson:"first_name"`
	LastName  	string `json:"last_name" bson:"last_name"`
	ID			primitive.ObjectID `json:"id" bson:"_id"`
}

// Creates a new user with request data and inserts into DB
/*
Cases:
-happy case
-invalid data (empty strings for fname/lname)
-missing fields
 */
func (uh UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createUser called")
	var newUser user

	// Validate and unmarshal to newUser
	rules := govalidator.MapData{
		"first_name": []string{"required"},
		"last_name": []string{"required"},
	}
	opts := govalidator.Options{
		Data:            &newUser,
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

	newUser.ID = primitive.NewObjectID()
	err := uh.Session.DB("admin-db").C(UserCollection).Insert(newUser); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := utils.HTTPResponseObject{
		Meta: 	nil,
		Data:	newUser,
	}
	json.NewEncoder(w).Encode(response)
}

// Returns a list of all users
/*
Cases:
-happy
-empty -> returns nil
 */
func (uh UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getAllUsers called")

	var results []user
	err := uh.Session.DB("admin-db").C(UserCollection).Find(nil).All(&results); if err != nil {
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

// Gets a single user by ID
/* Cases
-happy path
-not found
-bad syntax for id? or empty
 */
func (uh UserHandler) GetSingleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getSingleUser called")
	userID := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage: 	"Bad id syntax",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	var user user
	err = uh.Session.DB("admin-db").C(UserCollection).FindId(objId).One(&user); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id ID not found", // TODO figure out string interpolation
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errBody)
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

// Updates user by ID
/* Cases
-happy path
TODO figure out which things i can update and pass in
 */
func (uh UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: updateUser called")
	userID := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage: 	"Bad id syntax",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody)
		return
	}
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		// TODO custom error
		return
	}
	json.Unmarshal(reqBody, &updatedUser)
	updatedUser.ID = objId

	err = uh.Session.DB("admin-db").C(UserCollection).UpdateId(objId, updatedUser); if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(updatedUser)
	w.WriteHeader(http.StatusOK)
}

func (uh UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: deleteUser called")

	userID := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage: 	"Bad id syntax",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errBody)
		return
	}

	err = uh.Session.DB("admin-db").C(UserCollection).RemoveId(objId); if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // try again for not found
		// TODO custom error
	}
	w.WriteHeader(http.StatusNoContent)
}