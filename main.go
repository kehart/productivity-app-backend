package main

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/handlers"
	"github.com/productivity-app-backend/managers"
	"github.com/productivity-app-backend/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
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


func main() {

	initValidator()

	s := getSession()
	adminStore := utils.MongoDb{
		Session:        s,
		DbName:         utils.DbName,
	}
	userManager := 	managers.UserManagerImpl{Store: &adminStore}
	uh := handlers.UserHandler{ UserManager: &userManager }
	goalHandler := managers.GoalManagerImpl{Store: &adminStore}
	gh := handlers.GoalHandler{GoalManager:goalHandler}
	eventManager := managers.EventManagerImpl{Store: &adminStore}
	eh := handlers.EventHandler{EventManager: eventManager}

	router := mux.NewRouter().StrictSlash(true)

	// Test path
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test succeeded"))
	}).Methods("GET")

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
	router.HandleFunc("/events", eh.GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", eh.GetSingleEvent).Methods("GET")

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":" + port, router)) // create server
}