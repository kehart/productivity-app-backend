package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/managers"
	"github.com/productivity-app-backend/src/utils"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)
const (
	UserCollection = "users"
)

type UserHandler struct {
	UserManager *managers.UserManager
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
	var newUser utils.User

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

	err := uh.UserManager.CreateUser(&newUser); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
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

	var results []utils.User
	_, err := uh.UserManager.GetUsers(&results); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
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
	objId, err := formatObjectId(userID); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	var user utils.User
	_, errLong := uh.UserManager.GetSingleUser(&user, objId); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

// Updates user by ID; should be able to update first name and last name
/* Cases
-happy path (change both or one field)
-not found
 */
func (uh UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: updateUser called")

	userID := mux.Vars(r)["id"]
	objId, errLong := formatObjectId(userID); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		return
	}
	var updatedUser, existingUser utils.User

	reqBody, err := ioutil.ReadAll(r.Body); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage: 	"Invalid syntax",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}
	json.Unmarshal(reqBody, &updatedUser)
	updatedUser.ID = objId

	_, errLong = uh.UserManager.UpdateUser(&existingUser, &updatedUser); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		return
	}

	json.NewEncoder(w).Encode(existingUser)
	w.WriteHeader(http.StatusOK)
}

// cases
// -happy path
// not found
func (uh UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: deleteUser called")

	userID := mux.Vars(r)["id"]
	objId, err := formatObjectId(userID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	err = uh.UserManager.DeleteUser(objId); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func formatObjectId(userID string) (primitive.ObjectID, *utils.HTTPErrorLong) {
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Bad id syntax",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return objId, &fullErr
	}
	return objId, nil
}