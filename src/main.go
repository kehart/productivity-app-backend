package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type user struct {
	FirstName 	string `json:"first_name"`
	LastName  	string `json:"last_name"`
	ID			uuid.UUID `json:"id"`
}

type allUsers []user

var users = allUsers{
	{
		FirstName: "Kasia",
		LastName:  "Hart",
		ID: uuid.Must(uuid.NewV4()),
	},
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the  user only")
	}

	json.Unmarshal(reqBody, &newUser)
	newUser.ID = uuid.Must(uuid.NewV4())
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user data in order to update")
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i, singleUser := range users {
		if singleUser.ID.String() == userID {
			singleUser.FirstName = updatedUser.FirstName
			singleUser.LastName = updatedUser.LastName
			users = append(users[:i], singleUser)
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
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
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router)) // create server
}