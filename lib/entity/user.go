package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	CustomerID          primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	GovernmentID        string             `bson:"government_id" json:"government_id"`
	Email               string             `bson:"email" json:"email"`
	Alias               *string            `bson:"alias" json:"alias"`
	TemporaryAccessCode string             `bson:"temporary_access_code" json:"-"`
	Password            string             `bson:"password" json:"-"`
	CreatedBy           string             `bson:"created_by" json:"created_by"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           *time.Time         `bson:"updated_at" json:"updated_at"`
	UpdatedBy           *string            `bson:"updated_by" json:"updated_by"`
}
