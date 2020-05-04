package interfaces

import (
	"github.com/productivity-app-backend/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
)

type (
	// User Manager Abstraction
	IUserManager interface {
		CreateUser(newUser *models.User) *models.HTTPErrorLong
		GetUsers() (*[]models.User, *models.HTTPErrorLong)
		GetSingleUser(objId primitive.ObjectID) (*models.User, *models.HTTPErrorLong)
		UpdateUser(userId primitive.ObjectID, updatesToApply *models.User) (*models.User, *models.HTTPErrorLong)
		DeleteUser(objId primitive.ObjectID) *models.HTTPErrorLong
	}

	// Goal Manager Abstraction
	IGoalManager interface {
		CreateGoal(newGoal *models.Goal) (*models.Goal, *models.HTTPErrorLong)
		GetSingleGoal(objId primitive.ObjectID) (*models.Goal, *models.HTTPErrorLong)
		GetGoals(queryVals *url.Values) (*[]models.Goal, *models.HTTPErrorLong)
		DeleteGoal(objId primitive.ObjectID) *models.HTTPErrorLong
	}
	// Event Manager Abstraction
	IEventManager interface {
		CreateEvent(event *IEvent) (*IEvent, *models.HTTPErrorLong)
		GetSingleEvent(objId primitive.ObjectID) (*IEvent, *models.HTTPErrorLong)
		GetEvents(queryVals *url.Values) (*[]IEvent, *models.HTTPErrorLong)
	}
)