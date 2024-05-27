package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type DeletePlayer interface {
	Do(ctx context.Context, id string) error
}

type deletePlayer struct {
	DBWriter database.DBWriter
}

func NewDeletePlayer(dbWriter database.DBWriter) DeletePlayer {
	return &deletePlayer{
		DBWriter: dbWriter,
	}
}

func (uc *deletePlayer) Do(ctx context.Context, id string) error {
	return uc.DBWriter.DeletePlayer(context.Background(), id)
}
