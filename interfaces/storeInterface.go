package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	// db abstraction
	Store interface {
		Create(obj interface{}, collectionName string) error
		Delete(id primitive.ObjectID, collectionName string) error
		FindAll(collectionName string, dest interface{}, query ...*map[string]interface{}) error
		FindById(id primitive.ObjectID, collectionName string, dest interface{}) error
		Update(id primitive.ObjectID, obj interface{}, collectionName string) error
	}
)