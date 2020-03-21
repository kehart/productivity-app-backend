package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

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
