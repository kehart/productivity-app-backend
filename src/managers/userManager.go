package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

// Concrete implementation of UserManager interface
type UserManagerImpl struct {
	Store interfaces.Store
}

func (um UserManagerImpl) CreateUser(newUser *models.User) *models.HTTPErrorLong {
	log.Print(utils.InfoLog + "UserManager:CreateUser called")

	// Assign new ID to new user
	newUser.ID = primitive.NewObjectID()

	// Insert user into DB
	err := um.Store.Create(newUser, utils.UserCollection); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return &fullErr
	}
	return nil
}

func (um UserManagerImpl ) GetUsers() (*[]models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:GetUsers called")

	var results []models.User
	err := um.Store.FindAll(utils.UserCollection, &results); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error: 		errBody,
			StatusCode: http.StatusInternalServerError,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}

	//var users []models.User
	//for _, u := range results {
	//	fmt.Println(u)
	//	users = append(users, u.(models.User))
	//}
	return &results, nil
}

func (um UserManagerImpl ) GetSingleUser(objId primitive.ObjectID) (*models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:GetSingleUser called")

	//user, err := um.Store.FindById(objId, utils.UserCollection); if err != nil {
	var user models.User
	err := um.Store.FindById(objId, utils.UserCollection, &user); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("User with id %s not found", objId.String()),
		}
		fullErr := models.HTTPErrorLong {
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &user, nil
}

// updatedUser contains all the information for the update, including the ID of the user
func (um UserManagerImpl ) UpdateUser(userId primitive.ObjectID, updatesToApply *models.User) (*models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:UpdateUser called")

	// Read the current state of the user from the DB and place data into existingUser
	var existingUser2 models.User
	//obj, err := um.Store.FindById(userId, utils.UserCollection); if err != nil {
	err := um.Store.FindById2(userId, utils.UserCollection, &existingUser2); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("User with id %s not found", userId.String()),
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}

	//existingUser := obj.(*models.User)

	// Make changes to existing user based on updatesToApply data
	if len(updatesToApply.FirstName) > 0 {
		existingUser2.FirstName = updatesToApply.FirstName
	}
	if len(updatesToApply.LastName) > 0 {
		existingUser2.LastName = updatesToApply.LastName
	}

	err = um.Store.Update(userId, existingUser2, utils.UserCollection); if err != nil {
		errBody := models.HttpError{
			ErrorCode:    http.StatusText(http.StatusInternalServerError),
			ErrorMessage: utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &existingUser2, nil
}

func (um UserManagerImpl) DeleteUser(objId primitive.ObjectID) *models.HTTPErrorLong {
	log.Print(utils.InfoLog + "UserManager:DeleteUser called")

	err := um.Store.Delete(objId, utils.UserCollection); if err != nil {
		if err.Error() == "not found" {
			errBody := models.HttpError{
				ErrorCode:		http.StatusText(http.StatusNotFound),
				ErrorMessage: 	"ID not found",
			}
			fullErr := models.HTTPErrorLong{
				Error:      errBody,
				StatusCode: http.StatusNotFound,
			}
			log.Println(utils.ErrorLog + "Insert body here") // TODO ??
			return &fullErr
		}
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return &fullErr
	}
	return nil
}