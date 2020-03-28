package handlers

import (
	//"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/managers"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
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

// todo change this back to private
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost") // default mongo port

	// Check if connection error
	if err != nil {
		panic(err)
	}
	return s
}

func Router() *mux.Router {
	router := mux.NewRouter()
	um := managers.UserManager{ Session: getSession() }
	uh := UserHandler{ UserManager: &um }
	router.HandleFunc("/users", uh.GetAllUsers).Methods("GET")
	return router
}


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
	request, _ := http.NewRequest("GET", "/users", nil) // nil for body
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "OK response is expected")
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
func TestUserHandler_GetSingleUser(t *testing.T) {

}

// Updates user by ID; should be able to update first name and last name
/* Cases
-happy path (change both or one field)
-not found
*/
func TestUserHandler_UpdateUser(t *testing.T) {

}
