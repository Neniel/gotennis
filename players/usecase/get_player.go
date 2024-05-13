package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
)

type GetPlayerUsecase interface {
	Get(ctx context.Context, id string) (*entity.Player, error)
}

type getPlayerUsecase struct {
	DBReader database.DBReader
}

func NewGetPlayerUsecase(dbReader database.DBReader) GetPlayerUsecase {
	return &getPlayerUsecase{
		DBReader: dbReader,
	}
}

func (uc *getPlayerUsecase) Get(ctx context.Context, id string) (*entity.Player, error) {
	return uc.DBReader.GetPlayer(ctx, id)
}
