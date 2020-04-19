package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/productivity-app-backend/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeUserManager struct {
	mock.Mock
}

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
	//request, _ := http.NewRequest("GET", "/users", nil) // nil for body
	//response := httptest.NewRecorder()
	//Router().ServeHTTP(response, request)
	//assert.Equal(t, http.StatusOK, response.Code, "OK response is expected")
}

// cases
// -happy path
// not found
func TestUserHandler_DeleteUser(t *testing.T) {

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

func TestUserHandler_GetSingleUser(t *testing.T) {
	//e404 := utils.HTTPErrorLong{
	//	Error:      utils.HttpError{},
	//	StatusCode: http.StatusNotFound,
	//}
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
		//{
		//	shouldFail: true,
		//	user: nil,
		//	error: &e404,
		//},
	}

	url := "/users/" + id.Hex()

	for _, tc := range cases {
		fakeManager := new(fakeUserManager)
		handler := UserHandler{fakeManager}
		if tc.shouldFail {
			fakeManager.On("GetSingleUser", id).Return(tc.shouldFail, tc.error.StatusCode) // shouldFail, errorCode
		} else {
			fakeManager.On("GetSingleUser", id).Return(tc.shouldFail) // shouldFail, errorCode
		}


		r, _ := http.NewRequest(http.MethodGet, url, nil)

		//Normal testing stuff
		rr := httptest.NewRecorder()

		//http.HandleFunc("/users/{id}", handler.GetSingleUser)
		//http.HandlerFunc(handler.GetSingleUser).ServeHTTP(rr, r) // Calls GetSingleUSer

		sm := http.NewServeMux()
		sm.HandleFunc("/users/{id}", handler.GetSingleUser)
		sm.ServeHTTP(rr, r)

		returnCode := rr.Code
		returnObj, _ := ioutil.ReadAll(rr.Body)
		if tc.shouldFail {
			assert.Equal(t,returnCode, http.StatusNotFound)
			assert.NotNil(t, rr.Body)
			var err utils.HTTPErrorLong
			json.Unmarshal(returnObj, &err)
			assert.Equal(t, err, tc.error)
		} else {
			assert.Equal(t, returnCode, http.StatusOK)
			assert.NotNil(t, rr.Body)
			var u utils.User
			json.Unmarshal(returnObj, &u)
			assert.Equal(t, tc.user, u)
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

/*
UserManager Interface Implementation
 */

func (_m *fakeUserManager) CreateUser(newUser *utils.User) *utils.HTTPErrorLong {
	return nil
}

func (_m *fakeUserManager)  GetUsers() (*[]utils.User, *utils.HTTPErrorLong) {
	return nil, nil
}

// (shouldFail, errorCode)
func (_m *fakeUserManager)  GetSingleUser(objId primitive.ObjectID) (*utils.User, *utils.HTTPErrorLong) {
	ret := _m.Called(objId)

	fmt.Println("THIS WAS CALLED")

	shouldFail := ret.Bool(0); if shouldFail {
		fmt.Println("executred shouldDFail")
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
	fmt.Println("Returning this user", user)
	return &user, nil
}

func (_m *fakeUserManager)  UpdateUser(userId primitive.ObjectID, updatesToApply *utils.User) (*utils.User, *utils.HTTPErrorLong) {
	return nil, nil
}

func (_m *fakeUserManager)  DeleteUser(objId primitive.ObjectID) *utils.HTTPErrorLong {
	return nil
}