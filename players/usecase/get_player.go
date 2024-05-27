package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type GetPlayer interface {
	Do(ctx context.Context, id string) (*entity.Player, error)
}

type getPlayer struct {
	DBReader database.DBReader
}

func NewGetPlayer(dbReader database.DBReader) GetPlayer {
	return &getPlayer{
		DBReader: dbReader,
	}
}

func (uc *getPlayer) Do(ctx context.Context, id string) (*entity.Player, error) {
	return uc.DBReader.GetPlayer(ctx, id)
}
