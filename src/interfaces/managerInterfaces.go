package interfaces

import (
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
)

type (
	// User Manager Abstraction
	IUserManager interface {
		CreateUser(newUser *models.User) *utils.HTTPErrorLong
		GetUsers() (*[]models.User, *utils.HTTPErrorLong)
		GetSingleUser(objId primitive.ObjectID) (*models.User, *utils.HTTPErrorLong)
		UpdateUser(userId primitive.ObjectID, updatesToApply *models.User) (*models.User, *utils.HTTPErrorLong)
		DeleteUser(objId primitive.ObjectID) *utils.HTTPErrorLong
	}

	// Goal Manager Abstraction
	IGoalManager interface {
		CreateGoal(newGoal *models.Goal) (*models.Goal, *utils.HTTPErrorLong)
		GetSingleGoal(objId primitive.ObjectID) (*models.Goal, *utils.HTTPErrorLong)
		GetGoals(queryVals *url.Values) (*[]models.Goal, *utils.HTTPErrorLong)
		DeleteGoal(objId primitive.ObjectID) *utils.HTTPErrorLong
	}
	// Event Manager Abstraction
	IEventManager interface {
		CreateEvent(event *IEvent) (*IEvent, *utils.HTTPErrorLong)
		GetSingleEvent(objId primitive.ObjectID) (*IEvent, *utils.HTTPErrorLong)
		GetEvents(queryVals *url.Values) (*[]IEvent, *utils.HTTPErrorLong)
	}
)