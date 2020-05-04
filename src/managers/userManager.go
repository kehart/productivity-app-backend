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
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return &fullErr
	}
	return nil
}

func (um UserManagerImpl ) GetUsers() (*[]models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:GetUsers called")

	var results []models.User
	err := um.Store.FindAll(utils.UserCollection, &results); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}

	return &results, nil
}

func (um UserManagerImpl ) GetSingleUser(objId primitive.ObjectID) (*models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:GetSingleUser called")

	var user models.User
	err := um.Store.FindById(objId, utils.UserCollection, &user); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), fmt.Sprintf("User with id %s not found", objId.String()), http.StatusNotFound)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &user, nil
}

// updatedUser contains all the information for the update, including the ID of the user
func (um UserManagerImpl ) UpdateUser(userId primitive.ObjectID, updatesToApply *models.User) (*models.User, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "UserManager:UpdateUser called")

	// Read the current state of the user from the DB and place data into existingUser
	var existingUser models.User
	err := um.Store.FindById(userId, utils.UserCollection, &existingUser); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), fmt.Sprintf("User with id %s not found", userId.String()), http.StatusNotFound)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}

	// Make changes to existing user based on updatesToApply data
	if len(updatesToApply.FirstName) > 0 {
		existingUser.FirstName = updatesToApply.FirstName
	}
	if len(updatesToApply.LastName) > 0 {
		existingUser.LastName = updatesToApply.LastName
	}

	err = um.Store.Update(userId, existingUser, utils.UserCollection); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return nil, &fullErr
	}
	return &existingUser, nil
}

func (um UserManagerImpl) DeleteUser(objId primitive.ObjectID) *models.HTTPErrorLong {
	log.Print(utils.InfoLog + "UserManager:DeleteUser called")

	err := um.Store.Delete(objId, utils.UserCollection); if err != nil {
		if err.Error() == "not found" {
			log.Println(utils.ErrorLog + "Insert body here") // TODO ??
			fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), "ID not found", http.StatusNotFound)
			return &fullErr
		}
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		log.Println(utils.ErrorLog + "Insert body here") // TODO ??
		return &fullErr
	}
	return nil
}