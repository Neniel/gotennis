package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type UpdateTournamentRequest struct {
}

type UpdateTournament interface {
	Do(ctx context.Context, id string, request *UpdateTournamentRequest) (*entity.Tournament, error)
}

type updateTournament struct {
	DBWriter database.DBWriter
	DBReader database.DBReader
}

func NewUpdateTournament(dbWriter database.DBWriter, dbReader database.DBReader) UpdateTournament {
	return &updateTournament{
		DBWriter: dbWriter,
		DBReader: dbReader,
	}
}

func (u *updateTournament) Do(ctx context.Context, id string, request *UpdateTournamentRequest) (*entity.Tournament, error) {
	return nil, nil
}
