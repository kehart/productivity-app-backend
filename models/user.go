package models

import (
	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

/* User Type Definitions */
type User struct {
	FirstName 	string `json:"first_name" bson:"first_name" valid:"type(string)"`
	LastName  	string `json:"last_name" bson:"last_name" valid:"type(string)"`
	ID			primitive.ObjectID `json:"id" bson:"_id" valid:"-"`
}

// TODO change to add a list of errors like https://golang.hotexamples.com/examples/github.com.asaskevich.govalidator/-/StringLength/golang-stringlength-function-examples.html
func ValidateUser(user *User) *HTTPErrorLong {
	if !valid.StringLength(user.FirstName, "1", "30") {
		errBody := HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, first_name length must be between 1 and 30 characters",
		}
		fullErr := HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	} else if !valid.StringLength(user.LastName, "1", "30") {
		errBody := HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, last_name length must be between 1 and 30 characters",
		}
		fullErr := HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	}
	return nil
}
