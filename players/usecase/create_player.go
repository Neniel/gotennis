package usecase

import (
	"context"
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

	}

	if r.Birthdate.Before(time.Now().UTC()) {

	}

	return nil
}

type CreatePlayerUsecase interface {
	CreatePlayer(request *CreatePlayerRequest) (*entity.Player, error)
}

type createPlayerUsecase struct {
	DBWriter database.DBWriter
}

func NewCreateCategoryUsecase(app *app.App) CreatePlayerUsecase {
	return &createPlayerUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
	}
}

func (uc *createPlayerUsecase) CreatePlayer(request *CreatePlayerRequest) (*entity.Player, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	player := entity.NewPlayer(request.FirstName, request.MiddleName, request.LastName, request.Birthdate)

	return uc.DBWriter.AddPlayer(context.Background(), player)
}
