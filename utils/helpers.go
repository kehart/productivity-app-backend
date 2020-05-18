package utils

import (
	"encoding/json"
	"fmt"
	"github.com/productivity-app-backend/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"strconv"
)

// Used to validate and format ObjectIds for use in a Mongo DB
func FormatObjectId(userID string) (primitive.ObjectID, *models.HTTPErrorLong) {
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := models.HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Bad id syntax",
		}
		fullErr := models.HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return objId, &fullErr
	}
	return objId, nil
}

// Takes in url.Values (essentially a map of query KVPs) and returns a map
// of KVPs that can be used in a DB query (converted types where necessary)
func ParseQueryString(queryVals *url.Values) *map[string]interface{} {

	// Implements querying
	queryValsMap := map[string][]string(*queryVals)
	finalQueryVals := make(map[string]interface{}, len(queryValsMap))
	for k, v := range queryValsMap {
		if k == "target_value" {
			finalQueryVals[k], _ = strconv.Atoi(v[0]) // TODO handle error better here
		} else {
			finalQueryVals[k] = v[0] // take the first value from the []string
		}
	}
	return &finalQueryVals
}

// TODO better more standardized error handling
func ReturnWithError(w http.ResponseWriter, statusCode int, status string, errorMessage string) {
	errBody := models.HttpError{
		ErrorCode:		status,
		ErrorMessage:	errorMessage,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errBody)
}

func ReturnWithErrorLong(w http.ResponseWriter, err models.HTTPErrorLong) {
	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(err.Error)
}

func ReturnSuccess(w http.ResponseWriter, body interface{}, statusCode int) {
	if body != nil {
		json.NewEncoder(w).Encode(body)
	}
	w.WriteHeader(statusCode)
}


const (
	BadRequestMessage = "Bad request"
	InternalServerErrorMessage = "Server error"
)

func NotFoundErrorString(objectName string, id interface{}) string {
	return fmt.Sprintf("%s with id %s not found", objectName, id)
}

const (
	InfoLog = "INFO: "
	ErrorLog = "ERROR: "
)