package usecase

import (
	"context"
	"time"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type UpdateTournamentRequest struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Location  string           `json:"location"`
	StartDate time.Time        `json:"start_date"`
	EndDate   time.Time        `json:"end_date"`
	Category  *entity.Category `json:"category"`
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
	tournament, err := u.DBReader.GetTournament(ctx, id)
	if err != nil {
		return nil, err
	}

	tournament.Name = request.Name
	tournament.Location = request.Location
	tournament.StartDate = request.StartDate
	tournament.EndDate = request.EndDate
	tournament.Category = request.Category

	return u.DBWriter.UpdateTournament(ctx, tournament)
}
