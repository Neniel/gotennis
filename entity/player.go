package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	FirstName  string             `bson:"first_name" json:"first_name"`
	MiddleName string             `bson:"middle_name" json:"middle_name"`
	LastName   string             `bson:"last_name" json:"last_name"`
	Birthdate  time.Time          `bson:"birthdate" json:"birthdate"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
