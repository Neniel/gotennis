package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
)

type ListPlayers interface {
	Do(ctx context.Context) ([]entity.Player, error)
}

type listPlayers struct {
	DBReader database.DBReader
}

func NewListPlayers(dbReader database.DBReader) ListPlayers {
	return &listPlayers{
		DBReader: dbReader,
	}
}

func (uc *listPlayers) Do(ctx context.Context) ([]entity.Player, error) {
	players, err := uc.DBReader.GetPlayers(ctx)

	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}

	return players, nil
}
