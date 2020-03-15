package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

func (uc UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the  user only")
	}

	json.Unmarshal(reqBody, &newUser)
	uc.session.DB("admin-db").C("users").Insert(newUser)

	newUser.ID = primitive.NewObjectID()
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

func (uc UserController) getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func (uc UserController) getSingleUser(w http.ResponseWriter, r *http.Request) {
	// Some hacky conversions to get right type for id
	userID := []byte(mux.Vars(r)["id"])
	rv := bson.RawValue{
		Type: bsontype.ObjectID,
		Value: 	userID,
	}
	objId, success := rv.ObjectIDOK(); if !success {
		fmt.Println("error")
	}
	// Set up var that will hold requested user data
	var user user

	_, err := ioutil.ReadAll(r.Body); if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	})
	// search for user
	err = uc.session.DB("admin-db").C("users").FindId(objId).One(&user); if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

func (uc UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user data in order to update")
		fmt.Println("error")
	}
	json.Unmarshal(reqBody, &updatedUser)

	// search for user
	for i, singleUser := range users {
		if singleUser.ID.String() == userID {
			singleUser.FirstName = updatedUser.FirstName
			singleUser.LastName = updatedUser.LastName
			users = append(users[:i], singleUser)
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func (uc UserController) deleteUser(w http.ResponseWriter, r *http.Request) {
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