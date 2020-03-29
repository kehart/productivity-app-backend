package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserManager struct {
	//Session *mgo.Session
	Store utils.Store
}

func (um UserManager) CreateUser(newUser *utils.User) *utils.HTTPErrorLong {
	fmt.Println("LOG: Manager.CreateUser called")

	// Assign new ID to new user
	newUser.ID = primitive.NewObjectID()
	// Insert user into DB
	err := um.Store.Create(newUser); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return &fullErr
	}
	return nil
}

func (um UserManager) GetUsers(/*results *[]utils.User*/) (*[]utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.GetUsers called")

	var results *[]utils.User
	results, err := um.Store.FindAll(); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error: 		errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return results, nil
}

func (um UserManager) GetSingleUser(objId primitive.ObjectID) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.GetSingleUser called")

	//var user utils.User
	user, err := um.Store.FindById(objId); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("User with id %s not found", objId.String()),
		}
		fullErr := utils.HTTPErrorLong {
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return user, nil
}

// updatedUser contains all the information for the update, including the ID of the user
func (um UserManager) UpdateUser(userId primitive.ObjectID, updatesToApply *utils.User) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.UpdateUser called")

	// Read the current state of the user from the DB and place data into existingUser
	var existingUser *utils.User
	existingUser, err := um.Store.FindById(userId); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("User with id %s not found", userId.String()),
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}

	// Make changes to existing user based on updatesToApply data
	if len(updatesToApply.FirstName) > 0 {
		existingUser.FirstName = updatesToApply.FirstName
	}
	if len(updatesToApply.LastName) > 0 {
		existingUser.LastName = updatesToApply.LastName
	}

	existingUser, err = um.Store.Update(userId, existingUser); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:    http.StatusText(http.StatusInternalServerError),
			ErrorMessage: "Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return existingUser, nil
}

func (um UserManager) DeleteUser(objId primitive.ObjectID) *utils.HTTPErrorLong {
	fmt.Println("LOG: Manager.DeleteUser called")

	err := um.Store.Delete(objId); if err != nil {
		if err.Error() == "not found" {
			errBody := utils.HttpError{
				ErrorCode:		http.StatusText(http.StatusNotFound),
				ErrorMessage: 	"ID not found",
			}
			fullErr := utils.HTTPErrorLong{
				Error:      errBody,
				StatusCode: http.StatusNotFound,
			}
			return &fullErr
		}
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return &fullErr
	}
	return nil
}