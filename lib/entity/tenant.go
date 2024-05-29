package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tenant struct {
	ID                      primitive.ObjectID `bson:"_id" json:"id"`
	Name                    string             `bson:"name" json:"name"`
	PhoneNumber             string             `bson:"phone_number" json:"phone_number"`
	Email                   string             `bson:"email" json:"email"`
	Tier                    string             `bson:"tier" json:"tier"`
	MongoDBConnectionString string             `bson:"mongo_db_connection_string" json:"-"`
	DatabaseName            string             `bson:"database_name" json:"database_name"`
	CreatedBy               string             `bson:"created_by" json:"created_by"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               *time.Time         `bson:"updated_at" json:"updated_at"`
	UpdatedBy               *string            `bson:"updated_by" json:"updated_by"`
}
