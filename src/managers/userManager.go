package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Concrete implementation of UserManager interface
type UserManagerImpl struct {
	Store utils.Store
}

func (um UserManagerImpl) CreateUser(newUser *utils.User) *utils.HTTPErrorLong {
	fmt.Println("LOG: UserManager.CreateUser called")

	// Assign new ID to new user
	newUser.ID = primitive.NewObjectID()
	// Insert user into DB
	err := um.Store.Create(newUser, utils.UserCollection); if err != nil {
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

func (um UserManagerImpl ) GetUsers() (*[]utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: UserManager.GetUsers called")

//	var results *[]utils.User
	results, err := um.Store.FindAll(utils.UserCollection); if err != nil {
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
	return results.(*[]utils.User), nil
}

func (um UserManagerImpl ) GetSingleUser(objId primitive.ObjectID) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: UserManager.GetSingleUser called")

	user, err := um.Store.FindById(objId, utils.UserCollection); if err != nil {
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
	return user.(*utils.User), nil
}

// updatedUser contains all the information for the update, including the ID of the user
func (um UserManagerImpl ) UpdateUser(userId primitive.ObjectID, updatesToApply *utils.User) (*utils.User, *utils.HTTPErrorLong) {
	fmt.Println("LOG: UserManager.UpdateUser called")

	// Read the current state of the user from the DB and place data into existingUser
	obj, err := um.Store.FindById(userId, utils.UserCollection); if err != nil {
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

	existingUser := obj.(*utils.User)

	// Make changes to existing user based on updatesToApply data
	if len(updatesToApply.FirstName) > 0 {
		existingUser.FirstName = updatesToApply.FirstName
	}
	if len(updatesToApply.LastName) > 0 {
		existingUser.LastName = updatesToApply.LastName
	}

	updatedUser, err := um.Store.Update(userId, existingUser, utils.UserCollection); if err != nil {
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
	return updatedUser.(*utils.User), nil
}

func (um UserManagerImpl) DeleteUser(objId primitive.ObjectID) *utils.HTTPErrorLong {
	fmt.Println("LOG: UserManager.DeleteUser called")

	err := um.Store.Delete(objId, utils.UserCollection); if err != nil {
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
			StatusCode: http.StatusInternalServerError,
		}
		return &fullErr
	}
	return nil
}