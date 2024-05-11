package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
)

type ListPlayersUsecase interface {
	List(ctx context.Context) ([]entity.Player, error)
}

type listPlayersUsecase struct {
	DBReader database.DBReader
}

func NewListPlayersUsecase(app *app.App) ListPlayersUsecase {
	return &listPlayersUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *listPlayersUsecase) List(ctx context.Context) ([]entity.Player, error) {
	return uc.DBReader.GetPlayers(ctx)
}
