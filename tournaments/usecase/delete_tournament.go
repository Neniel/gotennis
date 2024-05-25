package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type DeleteTournament interface {
	Do(ctx context.Context, id string) error
}

type deleteTournament struct {
	DBWriter database.DBWriter
}

func NewDeleteTournament(dbWriter database.DBWriter) DeleteTournament {
	return &deleteTournament{
		DBWriter: dbWriter,
	}
}

func (u *deleteTournament) Do(ctx context.Context, id string) error {
	return nil
}
