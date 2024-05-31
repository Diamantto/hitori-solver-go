package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Person Model
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty" validate:"required,alpha"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty" validate:"required,alpha"`
}

// Matrix Model
type Matrix struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Coordinates []int32            `bson:"coordinates" json:"coordinates"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	User        Person             `bson:"user" json:"user"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
