package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	Name       string             `bson:"name" json:"name"`
	Location   string             `bson:"location" json:"location"`
	StartDate  time.Time          `bson:"start_date" json:"start_date"`
	EndDate    time.Time          `bson:"end_date" json:"end_date"`
	Category   *Category          `bson:"category" json:"category"`
}
