package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/src/handlers"
	"github.com/productivity-app-backend/src/managers"
	"github.com/productivity-app-backend/src/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	valid "github.com/asaskevich/govalidator"
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

func initValidator() {
	valid.SetFieldsRequiredByDefault(true)
}


func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	initValidator()

	s := getSession()
	userManager := 	managers.UserManager{Store: utils.MongoDb{s}}
	uh := handlers.UserHandler{ UserManager: &userManager }
	goalHandler := managers.GoalManager{Session: s }
	gh := handlers.GoalHandler{GoalManager:goalHandler}
	eventManager := managers.EventManager{Session:s}
	eh := handlers.EventHandler{SleepManager: eventManager}

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

	// Event Routing
	router.HandleFunc("/events", eh.CreateEvent).Methods("POST")
	router.HandleFunc("/events", eh.GetAllEventsByType).Methods("GET")
	router.HandleFunc("/events/{id}", eh.GetSingleEvent).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router)) // create server
}