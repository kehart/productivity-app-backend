package managers

import (
	"github.com/productivity-app-backend/src/interfaces"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/url"
)

type GoalManagerImpl struct {
	Store interfaces.Store
}

func (gm GoalManagerImpl) CreateGoal(newGoal *models.Goal) (*models.Goal, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "GoalManager:CreateGoal called")

	// Check the userId in newGoal exists
	var user models.User
	err := gm.Store.FindById(newGoal.UserId, utils.UserCollection, &user); if err != nil { // TODO verify this condition works
		errBody := models.HttpError{
			ErrorCode: 		http.StatusText(http.StatusBadRequest),
			ErrorMessage:	utils.NotFoundErrorString("User", newGoal.UserId.String()),
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return nil, &fullErr
	}

	// Insert goal into db
	newGoal.ID = primitive.NewObjectID()
	err = gm.Store.Create(newGoal, utils.GoalCollection); if err != nil {
		errBody := models.HttpError{
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}
	return newGoal, nil
}

func (gm GoalManagerImpl) GetSingleGoal(objId primitive.ObjectID) (*models.Goal, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "GoalManager:GetSingleGoal called")

	var goal models.Goal
	err := gm.Store.FindById(objId, utils.GoalCollection, &goal); if err != nil {
		errBody := models.HttpError {
			ErrorCode:		http.StatusText(http.StatusNotFound),
			ErrorMessage: 	utils.NotFoundErrorString("Goal", objId.String()),
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusNotFound,
		}
		return nil, &fullErr
	}
	return &goal, nil
}

func (gm GoalManagerImpl) GetGoals(queryVals *url.Values) (*[]models.Goal, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "GoalManager:GetGoals called")

	finalQueryVals := utils.ParseQueryString(queryVals)
	var goals []models.Goal
	err := gm.Store.FindAll(utils.GoalCollection, &goals, finalQueryVals) ; if err != nil {
		errBody := models.HttpError {
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return nil, &fullErr
	}

	return &goals, nil
}

func (gm GoalManagerImpl) DeleteGoal(objId primitive.ObjectID) *models.HTTPErrorLong {
	log.Print(utils.InfoLog + "GoalManager:DeleteGoal called")

	err := gm.Store.Delete(objId, utils.GoalCollection); if err != nil {
		if err.Error() == "not found" {
			errBody := models.HttpError{
				ErrorCode:		http.StatusText(http.StatusNotFound),
				ErrorMessage: 	utils.NotFoundErrorString("Goal", objId.String()),
			}
			fullErr := models.HTTPErrorLong{
				Error:      errBody,
				StatusCode: http.StatusNotFound,
			}
			return &fullErr
		}
		errBody := models.HttpError {
			ErrorCode:		http.StatusText(http.StatusInternalServerError),
			ErrorMessage: 	utils.InternalServerErrorMessage,
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusInternalServerError,
		}
		return &fullErr
	}
	return nil
}