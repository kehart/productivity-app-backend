package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"log"
)

const (
	DbName 			= "admin-db"
	UserCollection 	= "users"
	GoalCollection 	= "goals"
	EventCollection = "events"

)

/*
MongoDb struct and Store interface implementation
 */

type MongoDb struct {
	Session 		*mgo.Session
	DbName 			string
}

func (mdb MongoDb) Create(obj interface{}, collectionName string) error {
	log.Print(InfoLog + "MongoDb:Create called")
	return mdb.Session.DB(mdb.DbName).C(collectionName).Insert(obj)
}

func (mdb MongoDb) FindById(id primitive.ObjectID, collectionName string, dest interface{}) error {
	log.Print(InfoLog + "MongoDb:FindById called")
	err := mdb.Session.DB(mdb.DbName).C(collectionName).FindId(id).One(dest)
	return err
}

func (mdb MongoDb) FindAll(collectionName string, dest interface{}, query ...*map[string]interface{}) error {
	log.Print(InfoLog + "MongoDb:FindAll called")
	var err error
	if len(query) > 0 {
		err = mdb.Session.DB(mdb.DbName).C(collectionName).Find(query[0]).All(dest)
	} else {
		err = mdb.Session.DB(mdb.DbName).C(collectionName).Find(nil).All(dest)
	}
	fmt.Println(err)
	return  err
}

func (mdb MongoDb) Update(id primitive.ObjectID, obj interface{}, collectionName string) error {
	log.Print(InfoLog + "MongoDb:Update called")
	err := mdb.Session.DB(mdb.DbName).C(collectionName).UpdateId(id, obj)
	return err
}

func (mdb MongoDb) Delete(id primitive.ObjectID, collectionName string) error {
	log.Print(InfoLog + "MongoDb:Delete called")
	return  mdb.Session.DB(mdb.DbName).C(collectionName).RemoveId(id)
}

/*
RELATIONAL Db struct and Store interface implementation
 */

