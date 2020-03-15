package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
)

type UserController struct{
	session 	*mgo.Session
}


type user struct {
	FirstName 	string `json:"first_name" bson:"first_name"`
	LastName  	string `json:"last_name" bson:"last_name"`
	ID			primitive.ObjectID `json:"id" bson:"_id"`
}

type allUsers []user

var users = allUsers{
	{
		FirstName: "Kasia",
		LastName:  "Hart",
		ID: primitive.NewObjectID(),
	},
}

// session *mgo.session

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost") // default mongo port

	// Check if connection error
	if err != nil {
		panic(err)
	}
	return s
}

// Creates a new user with request data and inserts into DB
func (uc UserController) createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: createUser called")

	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the  user only")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &newUser)
	newUser.ID = primitive.NewObjectID()
	err = uc.session.DB("admin-db").C("users").Insert(newUser); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

// Returns a list of all users
func (uc UserController) getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getAllUsers called")

	var results []user
	err := uc.session.DB("admin-db").C("users").Find(nil).All(&results); if err != nil {
		// TODO: what should actually happen here?
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: what do i return if empty
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

// Gets a single user by ID
func (uc UserController) getSingleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: getSingleUser called")
	userID := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("error\n")
		return
	}
	// Set up var that will hold requested user data
	var user user

	_, err = ioutil.ReadAll(r.Body); if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// search for user
	err = uc.session.DB("admin-db").C("users").FindId(objId).One(&user); if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK) // TODO: superflous?
}

// Updates user by ID
func (uc UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: updateUser called")
	userID := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(userID); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("invalid id\n")
		return
	}
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user data in order to update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.Unmarshal(reqBody, &updatedUser)
	updatedUser.ID = objId

	err = uc.session.DB("admin-db").C("users").UpdateId(objId, updatedUser); if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(updatedUser)
	w.WriteHeader(http.StatusOK)
}

func (uc UserController) deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOG: deleteUser called")
	userID := mux.Vars(r)["id"]

	for i, singleUser := range users {
		if singleUser.ID.String() == userID {
			users = append(users[:i], users[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", userID)
		}
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	s := getSession()
	uc := UserController{session:s}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	router.HandleFunc("/users", uc.createUser).Methods("POST")
	router.HandleFunc("/users", uc.getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", uc.getSingleUser).Methods("GET")
	router.HandleFunc("/users/{id}", uc.updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", uc.deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router)) // create server
}