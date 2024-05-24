package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type DeletePlayerUsecase interface {
	Delete(ctx context.Context, id string) error
}

type deletePlayerUsecase struct {
	DBWriter database.DBWriter
}

func NewDeletePlayerUsecase(dbWriter database.DBWriter) DeletePlayerUsecase {
	return &deletePlayerUsecase{
		DBWriter: dbWriter,
	}
}

func (uc *deletePlayerUsecase) Delete(ctx context.Context, id string) error {
	return uc.DBWriter.DeletePlayer(context.Background(), id)
}
