package managers

import (
	"context"
	"github.com/productivity-app-backend/interfaces"
	"github.com/productivity-app-backend/models"
	"github.com/productivity-app-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusBadRequest), utils.NotFoundErrorString("User", newGoal.UserId.String()), http.StatusBadRequest)
		return nil, &fullErr
	}

	// Insert goal into db
	newGoal.ID = primitive.NewObjectID()
	err = gm.Store.Create(newGoal, utils.GoalCollection); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		return nil, &fullErr
	}
	return newGoal, nil
}

func (gm GoalManagerImpl) GetSingleGoal(objId primitive.ObjectID) (*models.Goal, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "GoalManager:GetSingleGoal called")

	var goal models.Goal
	err := gm.Store.FindById(objId, utils.GoalCollection, &goal); if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), utils.NotFoundErrorString("Goal", objId.String()), http.StatusNotFound)
		return nil, &fullErr
	}
	return &goal, nil
}

func (gm GoalManagerImpl) GetGoals(queryVals *url.Values) (*[]models.Goal, *models.HTTPErrorLong) {
	log.Print(utils.InfoLog + "GoalManager:GetGoals called")

	finalQueryVals := utils.ParseQueryString(queryVals)

	var goals []models.Goal
	decoder := func (cur *mongo.Cursor) error {

		for cur.Next(context.TODO()) {
			var goal models.Goal
			err := cur.Decode(&goal); if err != nil {
				return  err
			}
			goals = append(goals, goal)
		}

		err := cur.Err()
		return err
	}

	var bsonGoals []interface{}
	err := gm.Store.FindAll(utils.GoalCollection, bsonGoals, decoder, finalQueryVals) ; if err != nil {
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		return nil, &fullErr
	}

	return &goals, nil
}

func (gm GoalManagerImpl) DeleteGoal(objId primitive.ObjectID) *models.HTTPErrorLong {
	log.Print(utils.InfoLog + "GoalManager:DeleteGoal called")

	err := gm.Store.Delete(objId, utils.GoalCollection); if err != nil {
		if err.Error() == "not found" {
			fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusNotFound), utils.NotFoundErrorString("Goal", objId.String()), http.StatusNotFound)
			return &fullErr
		}
		fullErr := models.NewHTTPErrorLong(http.StatusText(http.StatusInternalServerError), utils.InternalServerErrorMessage, http.StatusInternalServerError)
		return &fullErr
	}
	return nil
}