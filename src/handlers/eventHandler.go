package handlers

import "go.mongodb.org/mongo-driver/bson/primitive"

// Should be abstract but for now just implement for sleep

type BaseEvent struct {
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id"`
	// dates
	// goal category
}
