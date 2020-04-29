package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	// db abstraction
	Store interface {
		Create(obj interface{}, collectionName string) error
		Delete(id primitive.ObjectID, collectionName string) error
		FindById(id primitive.ObjectID, collectionName string) (interface{}, error)
		FindAll(collectionName string, query ...*map[string]interface{}) (interface{}, error)
		Update(id primitive.ObjectID, obj interface{}, collectionName string) (interface{}, error)
	}
)