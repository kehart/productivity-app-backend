package utils

import (
	"context"
	"fmt"

	//"github.com/productivity-app-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

/* Second MongoDB */
type MongoDb2 struct {
	Session 		*mongo.Client
	DbName 			string
}


func (mdb MongoDb2) Create(obj interface{}, collectionName string) error {
	log.Print(InfoLog + "MongoDb2:Create called")
	_, err := mdb.Session.Database(mdb.DbName).Collection(collectionName).InsertOne(context.TODO(), obj)
	return err
}

func (mdb MongoDb2) FindById(id primitive.ObjectID, collectionName string, dest interface{}) error {
	log.Print(InfoLog + "MongoDb2:FindById called")
	filter := bson.D{{"_id", id}}
	err := mdb.Session.Database(mdb.DbName).Collection(collectionName).FindOne(context.TODO(), filter).Decode(dest)
	return err
}

func (mdb MongoDb2) FindAll(collectionName string, dest []interface{},
							decoder func(cur *mongo.Cursor) error,
							query ...*map[string]interface{}) error {

	// TODO remove dest
	log.Print(InfoLog + "MongoDb2:FindAll called")
	var err error
	var cur *mongo.Cursor
	findOptions := options.Find()

	if len(query) > 0 {
		//err = mdb.Session.Database(mdb.DbName).Collection(collectionName).Find(query[0]).All(dest)
		cur, err = mdb.Session.Database(mdb.DbName).Collection(collectionName).Find(context.TODO(), query[0], findOptions)
	} else {
		//err = mdb.Session.Database(mdb.DbName).Collection(collectionName).Find(nil).All(dest)
		cur, err = mdb.Session.Database(mdb.DbName).Collection(collectionName).Find(context.TODO(), bson.D{}, findOptions)
	}

	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())
	err = decoder(cur)
	return err
}

// Todo make this better https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
func (mdb MongoDb2) Update(id primitive.ObjectID, obj interface{}, collectionName string) error {
	log.Print(InfoLog + "MongoDb2:Update called")
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{
			"$set", obj,
		},
	}
	_, err := mdb.Session.Database(mdb.DbName).Collection(collectionName).UpdateOne(context.TODO(), filter, update) //(id, obj)
	return err
}

func (mdb MongoDb2) Delete(id primitive.ObjectID, collectionName string) error {
	log.Print(InfoLog + "MongoDb2:Delete called")
	filter := bson.D{{"_id", id}}
	_, err :=  mdb.Session.Database(mdb.DbName).Collection(collectionName).DeleteOne(context.TODO(), filter) //id)
	return err
}

