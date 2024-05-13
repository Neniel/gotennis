package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
)

type ListPlayersUsecase interface {
	List(ctx context.Context) ([]entity.Player, error)
}

type listPlayersUsecase struct {
	DBReader database.DBReader
}

func NewListPlayersUsecase(dbReader database.DBReader) ListPlayersUsecase {
	return &listPlayersUsecase{
		DBReader: dbReader,
	}
}

func (uc *listPlayersUsecase) List(ctx context.Context) ([]entity.Player, error) {
	return uc.DBReader.GetPlayers(ctx)
}
