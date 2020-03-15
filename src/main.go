package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/handlers"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost") // default mongo port

	// Check if connection error
	if err != nil {
		panic(err)
	}
	return s
}


func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	s := getSession()
	uh := handlers.UserHandler{Session: s}
	gh := handlers.GoalHandler{Session: s}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	// User Routing
	router.HandleFunc("/users", uh.CreateUser).Methods("POST")
	router.HandleFunc("/users", uh.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", uh.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", uh.UpdateUser).Methods("PUT") // TODO: change to patch
	router.HandleFunc("/users/{id}", uh.DeleteUser).Methods("DELETE")

	// Goal Routing
	router.HandleFunc("/goals", gh.CreateGoal).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router)) // create server
}