package entity

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	GovernmentID        string             `bson:"government_id" json:"government_id"`
	FirstName           string             `bson:"first_name" json:"first_name"`
	MiddleName          string             `bson:"middle_name" json:"middle_name"`
	LastName            string             `bson:"last_name" json:"last_name"`
	Birthdate           *time.Time         `bson:"birthdate" json:"birthdate"`
	PhoneNumber         string             `bson:"phone_number" json:"phone_number"`
	Email               string             `bson:"email" json:"email"`
	Alias               *string            `bson:"alias" json:"alias"`
	TemporaryAccessCode string             `bson:"temporary_access_code" json:"-"`
	Password            string             `bson:"password" json:"-"`
	Category            *Category          `bson:"category" json:"category"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           *time.Time         `bson:"updated_at" json:"updated_at"`
}

func NewPlayer(
	governmentID string,
	firstName string,
	middleName string,
	lastName string,
	birthDate *time.Time,
	phoneNumber string,
	email string,
	alias *string,
) *Player {
	return &Player{
		GovernmentID: governmentID,
		FirstName:    firstName,
		MiddleName:   middleName,
		LastName:     lastName,
		Birthdate:    birthDate,
		PhoneNumber:  phoneNumber,
		Email:        email,
		Alias:        alias,
		Category:     nil,
	}
}

func (p *Player) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Player) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
