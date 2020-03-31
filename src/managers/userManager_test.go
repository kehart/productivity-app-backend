package	managers

import (
	"errors"
	"github.com/productivity-app-backend/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

type createUserTest struct {
	user 		*utils.User
	error 		*utils.HTTPErrorLong
	shouldFail 	bool
}

type getUserTest struct {
	numUsers	int
	shouldFail	bool
	error 		*utils.HTTPErrorLong
}

func TestInsertUser(t *testing.T) {
	assert := assert.New(t)

	id := primitive.ObjectID{}
	u := utils.User{
		FirstName: "Bruce",
		LastName:  "Lee",
		ID:        id,
	}
	e := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusInternalServerError,
	}
	testCases := []createUserTest{
		{
			user: 		&u,
			error: 		nil,
			shouldFail: false,
		},
		{
			user: 		&u,
			error: 		&e,
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		db := new(fakeStore)
		manager := UserManager{Store:db}
		if tc.shouldFail {
			db.On("Create", &u).Return(errors.New("error"))
		} else {
			db.On("Create", &u).Return(nil)
		}

		var err *utils.HTTPErrorLong
		err = manager.CreateUser(&u) // calls create

		assert.Equal(tc.shouldFail, err != nil, "If the test case shouldFail, then the error must be nil") // todo change to expected/got with sprintf
		if tc.shouldFail {
			assert.Equal(tc.error.StatusCode, err.StatusCode) // only expecting internal server error
		} else {
			assert.NotEqual(tc.user.ID, id) // make sure the id was randomized
		}

	}
}


func TestGetEvents(t *testing.T) {
	assert := assert.New(t)

	e := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusInternalServerError,
	}
	testCases := []getUserTest{
		{
			numUsers: 0,
			shouldFail: true,
			error: &e,
		},
		{
			numUsers: 0,
			shouldFail: false,
			error: nil,
		},
		{
			numUsers: 5,
			shouldFail: false,
			error: nil,
		},
	}

	for _, tc := range testCases {
		db := new(fakeStore)
		manager := UserManager{Store:db}
		db.On("FindAll").Return(tc.shouldFail, tc.numUsers)

		var users *[]utils.User
		var err *utils.HTTPErrorLong
		users, err = manager.GetUsers() // calls FindAll

		assert.Equal(tc.shouldFail, err != nil, "If the test case shouldFail, then the error must be nil") // todo change to expected/got with sprintf
		if tc.shouldFail {
			assert.Equal(tc.error.StatusCode, err.StatusCode) // only expecting internal server error
		} else {
			assert.NotNil(users, "Users should not be nil")
			assert.Equal(len(*users), tc.numUsers, "Users should have the specified length")
		}
	}
}

func TestUserManager_GetSingleUser(t *testing.T) {
	assert := assert.New(t)

	e := utils.HTTPErrorLong{
		Error:      utils.HttpError{},
		StatusCode: http.StatusNotFound,
	}
	testCases := []getUserTest{
		{
			numUsers: 0,
			shouldFail: true,
			error: &e,
		},
		{
			numUsers: 1,
			shouldFail: false,
			error: nil,
		},
	}

	for _, tc := range testCases {
		db := new(fakeStore)
		manager := UserManager{Store:db}
		db.On("FindById").Return(tc.shouldFail)

		var user *utils.User
		var err *utils.HTTPErrorLong
		user, err = manager.GetSingleUser(primitive.ObjectID{}) // calls FindById

		assert.Equal(tc.shouldFail, err != nil, "If the test case shouldFail, then the error must be nil") // todo change to expected/got with sprintf
		if tc.shouldFail {
			assert.Equal(tc.error.StatusCode, err.StatusCode) // only expecting internal server error
		} else {
			assert.NotNil(user, "Users should not be nil")
		}

	}
}


/*
Mock Db and method implementations
 */

type fakeStore struct {
	mock.Mock
}

// either create returns nil as the error, or it returns a new error
func (_m *fakeStore) Create(user *utils.User) error {
	ret := _m.Called(user)

	var err error
	err, ok := ret.Get(0).(error); if ok { // casts value to type error
		return err
	} else if err != nil {
		panic("invalid type passed in")
	} else {
		return err // nil
	}
}
func (_m *fakeStore) Delete(Id primitive.ObjectID) error {
	//	ret := _m.Called(ID)
	//
	//	var r0 error
	//	if rf, ok := ret.Get(0).(func(int) error); ok {
	//		r0 = rf(ID)
	//	} else {
	//		r0 = ret.Error(0)
	//	}
	//
	//	return r0
	return nil
}

// takes in a boolean shouldFail
func (_m *fakeStore) FindById(id primitive.ObjectID) (*utils.User, error) {
	ret := _m.Called()

	shouldErr := ret.Bool(0); if shouldErr {
		return nil, errors.New("error")
	} else {
		user := utils.User{
			FirstName: "Bruce",
			LastName:  "Lee",
			ID:        primitive.ObjectID{},
		}
		return &user, nil
	}
}

// takes input of form (shouldErr, numUsers) and returns either
// an error or a slice of numUsers users
func (_m *fakeStore) FindAll() (*[]utils.User, error) {
	ret := _m.Called()

	shouldErr := ret.Bool(0); if shouldErr {
		return nil, errors.New("error")
	} else {
		numUsers := ret.Int(1)
		users := make([]utils.User, numUsers)
		return &users, nil
	}
}


func (_m *fakeStore) Update(id primitive.ObjectID, user *utils.User) (*utils.User, error) {
	return nil, nil
}

// TODO
/*
- change signatures in this module to match interface
- complete the mock functions
- write the test functions
 FOR THE HANDLERS
- create an interface for the manager
- create concrete implementation for the manager
- create mock implementation
- write handler functions
*/
