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

/*
MongoDb struct and Store interface implementation
 */

type MongoDb struct {
	Session 		*mgo.Session
	DbName 			string
}

func (mdb MongoDb) Create(obj interface{}, collectionName string) error {
	return mdb.Session.DB(mdb.DbName).C(collectionName).Insert(obj)
}

func (mdb MongoDb) FindById(id primitive.ObjectID, collectionName string, dest interface{}) error {
	err := mdb.Session.DB(mdb.DbName).C(collectionName).FindId(id).One(dest)
	return err
}

func (mdb MongoDb) FindAll(collectionName string, dest interface{}, query ...*map[string]interface{}) error {
	var err error
	if len(query) > 0 {
		err = mdb.Session.DB(mdb.DbName).C(collectionName).Find(query[0]).All(dest)
	} else{
		err = mdb.Session.DB(mdb.DbName).C(collectionName).Find(nil).All(dest)
	}
	return  err
}

func (mdb MongoDb) Update(id primitive.ObjectID, obj interface{}, collectionName string) error {
	err := mdb.Session.DB(mdb.DbName).C(collectionName).UpdateId(id, obj)
	return err
}

func (mdb MongoDb) Delete(id primitive.ObjectID, collectionName string) error {
	return  mdb.Session.DB(mdb.DbName).C(collectionName).RemoveId(id)
}

/*
RELATIONAL Db struct and Store interface implementation
 */

