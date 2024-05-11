package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type CreatePlayerRequest struct {
	FirstName  string           `json:"first_name"`
	MiddleName string           `json:"middle_name"`
	LastName   string           `json:"last_name"`
	Category   *entity.Category `json:"category"`
	Birthdate  time.Time        `json:"birthdate"`
}

func (r *CreatePlayerRequest) Validate() error {
	if r.FirstName == "" {
		return util.ErrPlayerFirstNameIsEmpty
	}

	if r.LastName == "" {
		return util.ErrPlayerLastNameIsEmpty
	}

	if r.Birthdate.IsZero() {
		return util.ErrPlayerBirthdateIsEmpty
	}

	if r.Birthdate.Before(time.Now().UTC()) {
		return util.ErrPlayerBirthdateIsFutureDate
	}

	return nil
}

type CreatePlayerUsecase interface {
	CreatePlayer(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error)
}

type createPlayerUsecase struct {
	DBWriter database.DBWriter
}

func NewCreatePlayerUsecase(app *app.App) CreatePlayerUsecase {
	return &createPlayerUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
	}
}

func (uc *createPlayerUsecase) CreatePlayer(ctx context.Context, request *CreatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(); err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when validating request: %w", err))
		return nil, err
	}

	newPlayer := entity.NewPlayer(request.FirstName, request.MiddleName, request.LastName, request.Birthdate)

	player, err := uc.DBWriter.AddPlayer(ctx, newPlayer)
	if err != nil {
		log.Println(fmt.Errorf("couldn't create player. Error when attempting add user to the database request: %w", err))
		return nil, err
	}
	return player, nil
}
