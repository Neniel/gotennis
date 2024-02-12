package entity

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

func (c *Category) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Category) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
