package managers

import (
	"fmt"
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
)

type GoalManagerImpl struct {
	Store interfaces.Store
}

func (gm GoalManagerImpl) CreateGoal(newGoal *models.Goal) (*models.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.CreateGoal called")

	// Check the userId in newGoal exists
	_, err := gm.Store.FindById(newGoal.UserId, utils.UserCollection); if err != nil { // TODO verify this condition works
		errBody := utils.HttpError{
			ErrorCode: 		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	"user with id user_id not found",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return nil, &fullErr
	}

	// Insert goal into db
	newGoal.ID = primitive.NewObjectID()
	err = gm.Store.Create(newGoal, utils.GoalCollection); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"internal server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return newGoal, nil
}

func (gm GoalManagerImpl) GetSingleGoal(objId primitive.ObjectID) (*models.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.GetSingleGoal called")

	goal, err := gm.Store.FindById(objId, utils.GoalCollection); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	fmt.Sprintf("Goal with id %s not found", objId.String()),
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return goal.(*models.Goal), nil
}

func (gm GoalManagerImpl) GetGoals(queryVals *url.Values) (*[]models.Goal, *utils.HTTPErrorLong) {
	fmt.Println("LOG: GoalManager.GetGoals called")

	finalQueryVals := utils.ParseQueryString(queryVals)
	results, err := gm.Store.FindAll(utils.GoalCollection, finalQueryVals); if err != nil {
		errBody := utils.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	"Server error",
		}
		fullErr := utils.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return results.(*[]models.Goal), nil
}

func (gm GoalManagerImpl) DeleteGoal(objId primitive.ObjectID) *utils.HTTPErrorLong {
	fmt.Println("LOG: GoalManager.DeleteGoal called")

	err := gm.Store.Delete(objId, utils.GoalCollection); if err != nil {
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