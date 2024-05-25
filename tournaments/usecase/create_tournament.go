package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type CreateTournamentRequest struct {
}

type CreateTournament interface {
	CreateTournament(ctx context.Context, request *CreateTournamentRequest) (*entity.Tournament, error)
}

type createTournament struct {
	DBWriter database.DBWriter
}

func NewCreateTournament(dbWriter database.DBWriter) CreateTournament {
	return &createTournament{
		DBWriter: dbWriter,
	}
}

func (u *createTournament) CreateTournament(ctx context.Context, request *CreateTournamentRequest) (*entity.Tournament, error) {
	return nil, nil
}
