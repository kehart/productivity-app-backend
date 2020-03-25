package utils

import (
	valid "github.com/asaskevich/govalidator"
	"net/http"
)

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
