package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	FirstName 	string `json:"first_name" bson:"first_name"`
	LastName  	string `json:"last_name" bson:"last_name"`
	ID			primitive.ObjectID `json:"id" bson:"_id"`
}
