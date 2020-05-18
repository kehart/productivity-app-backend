package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/models"
	"github.com/productivity-app-backend/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type getUserTest struct {
	shouldFail 		bool
	error			*utils.HTTPErrorLong
	user			*models.User
}

type bulkUserTest struct {
	users	[]models.User
}

/*
Pass in:
 */
func TestUserHandler_CreateUser(t *testing.T) {
	user := models.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        primitive.ObjectID{},
	}
	userFail := models.User{
		FirstName: 	"",
		LastName: 	"Lee",
		ID: 		primitive.ObjectID{},
	}
	e400 := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusBadRequest,
	}
	cases := []getUserTest{
		{
			shouldFail: false,
			user: &user,
			error: nil,
		},
		{
			shouldFail: true,
			user: &userFail,
			error: &e400,
		},
	}


	url := "/users"

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		if tc.shouldFail {
			fakeManager.On("CreateUser", tc.user).Return(tc.shouldFail, tc.error.StatusCode)
		} else {
			fakeManager.On("CreateUser", tc.user).Return(tc.shouldFail)
		}
		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(tc.user)

		r, _ := http.NewRequest(http.MethodPost, url, body)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users", handler.CreateUser).Methods(http.MethodPost)

		router.ServeHTTP(rr, r)

		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)

		if tc.shouldFail {
			assert.Equal(t, returnCode, http.StatusBadRequest)
			assert.NotNil(t, rr.Body)
		} else {
			assert.Equal(t, returnCode, http.StatusCreated)
			assert.NotNil(t, rr.Body)
			var user models.User
			json.Unmarshal(returnObj, &user)
			//assert.NotEqual(t, tc.user.ID, user.ID)
			assert.Equal(t, tc.user.FirstName, user.FirstName)
			assert.Equal(t, tc.user.LastName, user.LastName)
		}
	}
}


func TestUserHandler_GetAllUsers(t *testing.T) {
	cases := []bulkUserTest{
		{
			users: []models.User{ models.User{}, models.User{}, models.User{}, }, // TODO cleanup
		},
		{
			users: nil,
		},
	}

	url := "/users"

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		fakeManager.On("GetUsers").Return(tc.users)

		r, _ := http.NewRequest(http.MethodGet, url, nil)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users", handler.GetAllUsers).Methods(http.MethodGet)

		router.ServeHTTP(rr, r)


		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)

		assert.Equal(t, returnCode, http.StatusOK)
		assert.NotNil(t, rr.Body)
		var users []models.User
		json.Unmarshal(returnObj, &users)
		assert.Equal(t, tc.users, users)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	id := primitive.NewObjectID()
	user := models.User{	ID: id }
	e404 := utils.HTTPErrorLong {
		Error:      utils.HttpError{},
		StatusCode: http.StatusNotFound,
	}
	cases := []getUserTest{
		{
			shouldFail: false,
			user: &user,
			error: nil,
		},
		{
			shouldFail: true,
			user: &user,
			error: &e404,
		},
	}
	url := "/users/" + id.Hex()

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		if tc.shouldFail {
			fakeManager.On("DeleteUser", id).Return(tc.shouldFail, tc.error.StatusCode)
		} else {
			fakeManager.On("DeleteUser", id).Return(tc.shouldFail)
		}

		r, _ := http.NewRequest(http.MethodDelete, url, nil)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.DeleteUser).Methods(http.MethodDelete)

		router.ServeHTTP(rr, r)


		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)
		if tc.shouldFail {
			assert.Equal(t, returnCode, http.StatusNotFound)
			assert.NotNil(t, rr.Body)
			var err utils.HttpError
			json.Unmarshal(returnObj, &err)
			assert.Equal(t, tc.error.Error, err)
		} else {
			assert.Equal(t, returnCode, http.StatusNoContent)
		}
	}
}



func TestUserHandler_GetSingleUser(t *testing.T) {
	e404 := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusNotFound,
	}
	id := primitive.NewObjectID()
	user := models.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        id,
	}
	cases := []getUserTest{
		{
			shouldFail:false,
			user: &user,
			error: nil,
		},
		{
			shouldFail: true,
			user: nil,
			error: &e404,
		},
	}

	url := "/users/" + id.Hex()

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		if tc.shouldFail {
			fakeManager.On("GetSingleUser", id).Return(tc.shouldFail, tc.error.StatusCode)
		} else {
			fakeManager.On("GetSingleUser", id).Return(tc.shouldFail)
		}

		r, _ := http.NewRequest(http.MethodGet, url, nil)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.GetSingleUser).Methods("GET")

		router.ServeHTTP(rr, r)


		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)
		if tc.shouldFail {
			assert.Equal(t, returnCode, http.StatusNotFound)
			assert.NotNil(t, rr.Body)
			var err utils.HttpError
			json.Unmarshal(returnObj, &err)
			assert.Equal(t, tc.error.Error, err)
		} else {
			assert.Equal(t, returnCode, http.StatusOK)
			assert.NotNil(t, rr.Body)
			var u models.User
			json.Unmarshal(returnObj, &u)
			assert.Equal(t, *tc.user, u)
		}
	}
}

// Updates user by ID; should be able to update first name and last name
/* Cases
-happy path (change both or one field)
-not found
*/
func TestUserHandler_UpdateUser(t *testing.T) {
	// 404, 400, 200
	id := primitive.ObjectID{}
	userSuccess := models.User{
		FirstName: "Jackie",
		ID:       	id,
	}
	userFailure := models.User{
		FirstName: "",
		ID:        id,
	}
	e400 := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusBadRequest,
	}
	e404 := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusNotFound,
	}
	cases := []getUserTest{
		{
			shouldFail: false,
			user: &userSuccess,
			error: nil,
		},
		{
			shouldFail: true,
			user: &userFailure,
			error: &e400,
		},
		{
			shouldFail: true,
			user: &userFailure,
			error: &e404,
		},
	}

	url := "/users/" + id.Hex()

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		if tc.shouldFail {
			fakeManager.On("UpdateUser", id, tc.user).Return(tc.shouldFail, tc.error.StatusCode)
		} else {
			fakeManager.On("UpdateUser", id, tc.user).Return(tc.shouldFail)
		}

		body := new(bytes.Buffer)
		json.NewEncoder(body).Encode(tc.user)
		r, _ := http.NewRequest(http.MethodPatch, url, body)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.UpdateUser).Methods(http.MethodPatch)

		router.ServeHTTP(rr, r)


		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)
		if tc.shouldFail {
			assert.Equal(t, returnCode, tc.error.StatusCode)
			assert.NotNil(t, rr.Body)
			var err utils.HttpError
			json.Unmarshal(returnObj, &err)
			assert.Equal(t, tc.error.Error, err)
		} else {
			assert.Equal(t, returnCode, http.StatusOK)
			assert.NotNil(t, rr.Body)
			var u models.User
			json.Unmarshal(returnObj, &u)
			assert.Equal(t, tc.user.FirstName, u.FirstName)
			assert.Equal(t, tc.user.ID, u.ID)
		}
	}

}

// Implements IUserManager for use
// in test environment
type fakeUserManager struct {
	mock.Mock
}

/*
UserManager Interface Implementation
 */

func (_m *fakeUserManager) CreateUser(newUser *models.User) *utils.HTTPErrorLong {
	ret := _m.Called(newUser)

	shouldFail := ret.Bool(0); if shouldFail {
		errorCode := ret.Int(1)
		err := utils.HTTPErrorLong{
			Error:      utils.HttpError{},
			StatusCode: errorCode,
		}
		return &err
	}
	return nil
}

func (_m *fakeUserManager) GetUsers() (*[]models.User, *utils.HTTPErrorLong) {
	ret := _m.Called()

	users := ret.Get(0).([]models.User)
	return &users, nil
}

// (shouldFail, errorCode)
func (_m *fakeUserManager) GetSingleUser(objId primitive.ObjectID) (*models.User, *utils.HTTPErrorLong) {
	ret := _m.Called(objId)

	shouldFail := ret.Bool(0); if shouldFail {
		errorCode := ret.Int(1)
		err := utils.HTTPErrorLong{
			Error:      utils.HttpError{},
			StatusCode: errorCode,
		}
		return nil, &err
	}
	user := models.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        objId,
	}
	return &user, nil
}

// shouldFail, statusCode
func (_m *fakeUserManager) UpdateUser(userId primitive.ObjectID, updatesToApply *models.User) (*models.User, *utils.HTTPErrorLong) {
	ret := _m.Called(userId, updatesToApply)

	shouldFail := ret.Bool(0); if shouldFail {
		errorCode := ret.Int(1)
		err := utils.HTTPErrorLong{
			Error:      utils.HttpError{},
			StatusCode: errorCode,
		}
		return nil, &err
	}
	user := models.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        userId,
	}

	// Make changes to existing user based on updatesToApply data
	if len(updatesToApply.FirstName) > 0 {
		user.FirstName = updatesToApply.FirstName
	}
	if len(updatesToApply.LastName) > 0 {
		user.LastName = updatesToApply.LastName
	}

	return &user, nil
}

func (_m *fakeUserManager)  DeleteUser(objId primitive.ObjectID) *utils.HTTPErrorLong {
	ret := _m.Called(objId)

	shouldFail := ret.Bool(0); if shouldFail {
		errorCode := ret.Int(1)
		err := utils.HTTPErrorLong{
			Error:      utils.HttpError{},
			StatusCode: errorCode,
		}
		return &err
	}
	return nil
}