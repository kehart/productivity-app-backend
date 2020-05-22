package main

import (
	"context"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/productivity-app-backend/handlers"
	"github.com/productivity-app-backend/managers"
	"github.com/productivity-app-backend/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getSession() *mgo.Session {
	//// Connect to our local mongo
	//username := "rw-user"
	//pw := "5RxxooLgbikO26iA"
	//// "mongodb+srv://"
	/////test?retryWrites=true&w=majority
	//connString := "mongodb+srv://" + username + ":" + pw + "@cluster0-2taqh.mongodb.net/test"


	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//_, err := mongo.Connect(ctx, options.Client().ApplyURI(
	//	connString,
	//))
	//if err != nil { log.Fatal(err) }
	//fmt.Println(client)
	//fmt.Println("no erroe???")


	//tlsConfig := &tls.Config{}
	//
	//dialInfo := &mgo.DialInfo{
	//	Addrs: []string{"cluster0-2taqh.mongodb.net"},
	//	Database: "test",
	//	Username: username,
	//	Password: pw,
	//}
	//dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	//	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	//	return conn, err
	//}
	//session, err := mgo.DialWithInfo(dialInfo); if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(session)







	//
	//s, err := mgo.DialWithTimeout(connString, 30*time.Second) // default mongo port
	//log.Println(connString)
	//
	//// Check if connection error
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("returning from getsession")
	//return s


	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost") // default mongo port

	// Check if connection error
	if err != nil {
		panic(err)
	}
	return s
}

func getSession2() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func initValidator() {
	valid.SetFieldsRequiredByDefault(true)
}


func main() {

	initValidator()

	//s := getSession()
	s2 := getSession2()
	adminStore2 := utils.MongoDb2{
		Session: s2,
		DbName:  utils.DbName,
	}
	//adminStore := utils.MongoDb{
	//	Session:        s,
	//	DbName:         utils.DbName,
	//}
	userManager := 	managers.UserManagerImpl{Store: &adminStore2}
	uh := handlers.UserHandler{ UserManager: &userManager }
	goalHandler := managers.GoalManagerImpl{Store: &adminStore2}
	gh := handlers.GoalHandler{GoalManager:goalHandler}
	eventManager := managers.EventManagerImpl{Store: &adminStore2}
	eh := handlers.EventHandler{EventManager: eventManager}

	router := mux.NewRouter().StrictSlash(true)

	// Test path
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test succeeded"))
	}).Methods("GET")

	//// User Routing
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

	port := os.Getenv("PORT"); if port == "" {
		log.Fatal("No PORT specified")
	}

	log.Fatal(http.ListenAndServe(":" + port, router)) // create server
}