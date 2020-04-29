package handlers

import (
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"io/ioutil"
	"log"
	"net/http"
)

/*
A module for handling HTTP requests to the User API. Supports:
- Create
- Read Single and Read All
- Update
- Delete
 */

type UserHandler struct {
	UserManager interfaces.IUserManager
}

// Handles request for POST /users
// Takes in JSON user object
func (uh UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "UserHandler:CreateUser called")
	var newUser models.User

	reqBody, genErr := ioutil.ReadAll(r.Body); if genErr != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	utils.BadRequestMessage,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		log.Println(utils.ErrorLog + "Unable to read request body")
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
			log.Println(utils.ErrorLog + "Request body data invalid")
			return
	}
	err := models.ValidateUser(&newUser); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Request body data invalid") // TODO ??
		return
	}

	err = uh.UserManager.CreateUser(&newUser); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Handles request for GET /users
// Currently, no querying supported
func (uh UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "UserHandler:GetAllUsers called")

	var results *[]models.User
	results, err := uh.UserManager.GetUsers(); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

// Handles request for GET /users/{id}
func (uh UserHandler) GetSingleUser(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "UserHandler:GetSingleUser called")

	userID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(userID); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	user, errLong := uh.UserManager.GetSingleUser(objId); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

// Handles request for PATCH /users/{id}
// Takes in JSON user patch object
func (uh UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "UserHandler:UpdateUser called")

	// Extract and format id from URL
	userId := mux.Vars(r)["id"]
	objId, errLong := utils.FormatObjectId(userId); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	var updatesToApply models.User
	reqBody, err := ioutil.ReadAll(r.Body); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage: 	"Invalid syntax",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	json.Unmarshal(reqBody, &updatesToApply)

	// Validate fields in JSON object
	_, genErr := valid.ValidateStruct(&updatesToApply) ; if genErr != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	genErr,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errBody)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	updatedUser, errLong := uh.UserManager.UpdateUser(objId, &updatesToApply); if errLong != nil {
		w.WriteHeader(errLong.StatusCode)
		json.NewEncoder(w).Encode(errLong.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
	w.WriteHeader(http.StatusOK)
}

// Handles DELETE /users/{id}
func (uh UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Print(utils.InfoLog + "UserHandler:UpdateUser called")

	userID := mux.Vars(r)["id"]
	objId, err := utils.FormatObjectId(userID);  if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}

	err = uh.UserManager.DeleteUser(objId); if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(err.Error)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
