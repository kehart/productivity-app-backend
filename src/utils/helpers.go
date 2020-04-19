package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"strconv"
)

// Used to validate and format ObjectIds for use in a Mongo DB
func FormatObjectId(userID string) (primitive.ObjectID, *HTTPErrorLong) {
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		errBody := HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Bad id syntax",
		}
		fullErr := HTTPErrorLong{
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
