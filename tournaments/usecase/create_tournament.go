package usecase

import (
	"context"
	"time"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
)

type CreateTournamentRequest struct {
	Name      string           `json:"name"`
	Location  string           `json:"location"`
	StartDate time.Time        `json:"start_date"`
	EndDate   time.Time        `json:"end_date"`
	Category  *entity.Category `json:"category"`
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
	tournament := &entity.Tournament{
		Name:      request.Name,
		Location:  request.Location,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
		Category:  request.Category,
	}

	return u.DBWriter.AddTournament(ctx, tournament)
}
