package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Creates a new user with request data and inserts into DB
/*
Cases:
-happy case
-invalid data (empty strings for fname/lname)
-missing fields
*/



//func Router() *mux.Router {
//	router := mux.NewRouter()
//	um := managers.UserManager{ Session: getSession() }
//	uh := UserHandler{ UserManager: &um }
//	router.HandleFunc("/users", uh.GetAllUsers).Methods("GET")
//	return router
//}


/*
Pass in:
 */
//func TestUserHandler_CreateUser(t *testing.T) {
//	type test struct {
//		data []int
//		answer int
//	}
//	tests := []test{
//		{
//			data:   []int{1, 2, 3},
//			answer: 4,
//		},
//	}
//
//	for _, v := range tests {
//		// run the function
//		// check for err (expected, got)
//		fmt.Println(v)
//	}
//}

// Returns a list of all users
/*
Cases:
-happy
-empty -> returns nil
*/
func TestUserHandler_GetAllUsers(t *testing.T) {
	cases := []bulkUserTest{
		{
			users: []utils.User{ utils.User{}, utils.User{}, utils.User{}, },
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
		var users []utils.User
		json.Unmarshal(returnObj, &users)
		assert.Equal(t, tc.users, users)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	id := primitive.NewObjectID()
	user := utils.User{	ID: id }
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

// Gets a single user by ID
/* Cases
-happy path
-not found
-bad syntax for id? or empty
*/

type getUserTest struct {
	shouldFail 		bool
	error			*utils.HTTPErrorLong
	user			*utils.User
}

type bulkUserTest struct {
	users	[]utils.User
}

func TestUserHandler_GetSingleUser(t *testing.T) {
	e404 := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusNotFound,
	}
	id := primitive.NewObjectID()
	user := utils.User{
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
			var u utils.User
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

}

// Implements IUserManager for use
// in test environment
type fakeUserManager struct {
	mock.Mock
}

/*
UserManager Interface Implementation
 */

func (_m *fakeUserManager) CreateUser(newUser *utils.User) *utils.HTTPErrorLong {
	return nil
}

func (_m *fakeUserManager) GetUsers() (*[]utils.User, *utils.HTTPErrorLong) {
	ret := _m.Called()

	users := ret.Get(0).([]utils.User)
	return &users, nil
}

// (shouldFail, errorCode)
func (_m *fakeUserManager) GetSingleUser(objId primitive.ObjectID) (*utils.User, *utils.HTTPErrorLong) {
	ret := _m.Called(objId)

	shouldFail := ret.Bool(0); if shouldFail {
		errorCode := ret.Int(1)
		err := utils.HTTPErrorLong{
			Error:      utils.HttpError{},
			StatusCode: errorCode,
		}
		return nil, &err
	}
	user := utils.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        objId,
	}
	return &user, nil
}

func (_m *fakeUserManager)  UpdateUser(userId primitive.ObjectID, updatesToApply *utils.User) (*utils.User, *utils.HTTPErrorLong) {
	return nil, nil
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