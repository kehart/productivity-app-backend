package utils

import (
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"net/http"
)

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

// TODO change to have this as a receiver?
func ValidateGoal(goal *Goal) *HTTPErrorLong {
	if !goal.GoalCategory.isValid(goal.GoalType) {
		errBody := HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, goal_category and goal_name should be a valid pair",
		}
		fullErr := HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	} else if !goal.GoalType.isValid(goal.TargetValue) {
		errBody := HttpError{
			ErrorCode:    http.StatusText(http.StatusBadRequest),
			ErrorMessage: "Error, the type of target_value does not match goal_name",
		}
		fullErr := HTTPErrorLong{
			Error:      errBody,
			StatusCode: http.StatusBadRequest,
		}
		return &fullErr
	}
	return nil
}

// Validate goal type with respect to goal category
func (gc GoalCategory) isValid(gt GoalType) bool {
	switch gc {
	case Sleep:
		if gt == HoursSlept {
			return true
		}
		return false
	}
	return false
}

// Validate target value with respect to goal type
func (gt GoalType) isValid(target interface{}) bool {
	switch gt {
	case HoursSlept:
		tType := fmt.Sprintf("%T", target)
		if tType =="int" || tType == "float64" {
			return true
		}
		fmt.Println(tType)
		return false
	}
	return false
}

func (f Feeling) isValid() bool {
	switch f {
	case Sad, Happy, Tired, Anxious, Refreshed, Excited:
		return true
	}
	return false
}

/*
Notes for writing tests:
- file must end in _test.go
- put file in same package as the one being tested
- be in a func with signature func TestXxx(*testing.T)
- run test with go test
 */

