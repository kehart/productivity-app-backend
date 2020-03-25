package handlers

import (
	"encoding/json"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/managers"
	"github.com/productivity-app-backend/src/utils"
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
	}

	err = uh.UserManager.CreateUser(&newUser); if err != nil {
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
	objId, err := utils.FormatObjectId(userID); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		return
	}

	user, errLong := uh.UserManager.GetSingleUser(objId); if errLong != nil {
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
	objId, errLong := utils.FormatObjectId(userID); if errLong != nil {
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
	_, genErr := valid.ValidateStruct(&updatedUser) ; if genErr != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	genErr,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		return
	}
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
	objId, err := utils.FormatObjectId(userID);  if err != nil {
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
