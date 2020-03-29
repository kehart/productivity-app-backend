package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
)

const (
	DbName 			= "admin-db"
	UserCollection 	= "users"
	GoalCollection 	= "goals"
	EventCollection = "events"

)

type MongoDb struct {
	Session *mgo.Session
}

// Implement the Store interface

func (mdb MongoDb) Create(user *User) error {
	return  mdb.Session.DB(DbName).C(UserCollection).Insert(user)
}

func (mdb MongoDb) FindById(id primitive.ObjectID) (*User, error) {
	var user User
	err := mdb.Session.DB(DbName).C(UserCollection).FindId(id).One(&user)
	return &user, err
}

func (mdb MongoDb) FindAll() (*[]User, error) {
	var users []User
	err := mdb.Session.DB(DbName).C(UserCollection).Find(nil).All(&users)
	return &users, err
}

func (mdb MongoDb) Update(id primitive.ObjectID, user *User) (*User, error) {
	err := mdb.Session.DB(DbName).C(UserCollection).UpdateId(id, user)
	return user, err
}

func (mdb MongoDb) Delete(id primitive.ObjectID) error {
	return  mdb.Session.DB(DbName).C(UserCollection).RemoveId(id)
}

