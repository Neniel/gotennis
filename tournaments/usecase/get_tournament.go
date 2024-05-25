package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type GetTournament interface {
	Do(ctx context.Context, id string) (*entity.Tournament, error)
}

type getTournament struct {
	DBReader database.DBReader
}

func NewGetTournament(dbReader database.DBReader) GetTournament {
	return &getTournament{
		DBReader: dbReader,
	}
}

func (u *getTournament) Do(tx context.Context, id string) (*entity.Tournament, error) {
	return nil, nil
}
