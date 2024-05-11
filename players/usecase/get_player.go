package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
)

type GetPlayerUsecase interface {
	Get(ctx context.Context, id string) (*entity.Player, error)
}

type getPlayerUsecase struct {
	DBReader database.DBReader
}

func NewGetPlayerUsecase(app *app.App) GetPlayerUsecase {
	return &getPlayerUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *getPlayerUsecase) Get(ctx context.Context, id string) (*entity.Player, error) {
	return uc.DBReader.GetPlayer(ctx, id)
}
