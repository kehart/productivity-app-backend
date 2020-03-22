package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/handlers"
	"github.com/productivity-app-backend/src/managers"
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
	userManager := 	managers.UserManager{Session: s}
	uh := handlers.UserHandler{ UserManager: &userManager }
	goalHandler := managers.GoalManager{Session: s }
	gh := handlers.GoalHandler{GoalManager:goalHandler}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	// User Routing
	router.HandleFunc("/users", uh.CreateUser).Methods("POST")
	router.HandleFunc("/users", uh.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", uh.GetSingleUser).Methods("GET")
	router.HandleFunc("/users/{id}", uh.UpdateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", uh.DeleteUser).Methods("DELETE")

	// Goal Routing
	router.HandleFunc("/goals", gh.CreateGoal).Methods("POST")
	router.HandleFunc("/goals/{id}", gh.GetSingleGoal).Methods("GET")
	router.HandleFunc("/goals", gh.GetGoals).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router)) // create server
}