package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"net/http"
)

type UserManager struct {
	Session *mgo.Session
}

func (um UserManager) CreateUser(newUser *utils.User) *utils.HTTPErrorLong {
	fmt.Println("LOG: Manager.CreateUser called")

	// Assign new ID to new user
	newUser.ID = primitive.NewObjectID()
	// Insert user into DB
	err := um.Session.DB(utils.DbName).C(utils.UserCollection).Insert(newUser); if err != nil {
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

func (um UserManager) GetUsers(results *[]utils.User) (*[]utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.GetUsers called")

	err := um.Session.DB(utils.DbName).C(utils.UserCollection).Find(nil).All(results); if err != nil {
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

func (um UserManager) GetSingleUser(user *utils.User, objId primitive.ObjectID) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.GetSingleUser called")

	err := um.Session.DB(utils.DbName).C(utils.UserCollection).FindId(objId).One(&user); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id ID not found", // TODO figure out string interpolation
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
func (um UserManager) UpdateUser(existingUser *utils.User, updatedUser *utils.User) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: Manager.UpdateUser called")

	// Read the current state of the user from the DB and place data into existingUser
	err := um.Session.DB(utils.DbName).C(utils.UserCollection).FindId(updatedUser.ID).One(existingUser); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	"User with id ID not found", // TODO figure out string interpolation
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, & fullErr
	}

	// Make changes to existing user based on updateUser data
	if len(updatedUser.FirstName) > 0 {
		existingUser.FirstName = updatedUser.FirstName
	}
	if len(updatedUser.LastName) > 0 {
		existingUser.LastName = updatedUser.LastName
	}
	existingUser.ID = updatedUser.ID

	err = um.Session.DB(utils.DbName).C(utils.UserCollection).UpdateId(existingUser.ID, existingUser); if err != nil {
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

	err := um.Session.DB(utils.DbName).C(utils.UserCollection).RemoveId(objId); if err != nil {
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