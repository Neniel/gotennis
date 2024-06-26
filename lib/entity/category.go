package entity

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at" json:"updated_at"`
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}

func (c *Category) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Category) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
