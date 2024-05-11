package entity

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	FirstName  string             `bson:"first_name" json:"first_name"`
	MiddleName string             `bson:"middle_name" json:"middle_name"`
	LastName   string             `bson:"last_name" json:"last_name"`
	Birthdate  time.Time          `bson:"birthdate" json:"birthdate"`
	Category   *Category          `bson:"category" json:"category"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

func NewPlayer(firstName string, middleName string, lastName string, birthDate time.Time) *Player {
	return &Player{
		ID:         primitive.NewObjectID(),
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Birthdate:  birthDate,
		Category:   nil,
		CreatedAt:  time.Now().UTC(),
	}
}

func (p *Player) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Player) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
