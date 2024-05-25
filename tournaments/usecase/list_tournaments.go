package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type ListTournaments interface {
	Do(ctx context.Context) ([]entity.Tournament, error)
}

type listTournaments struct {
	DBReader database.DBReader
}

func NewListTournaments(dbReader database.DBReader) ListTournaments {
	return &listTournaments{
		DBReader: dbReader,
	}
}

func (u *listTournaments) Do(ctx context.Context) ([]entity.Tournament, error) {
	return nil
}
