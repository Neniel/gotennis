package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
)

type DeletePlayerUsecase interface {
	Delete(ctx context.Context, id string) error
}

type deletePlayerUsecase struct {
	DBWriter database.DBWriter
}

func NewDeletePlayerUsecase(app *app.App) DeletePlayerUsecase {
	return &deletePlayerUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
	}
}

func (uc *deletePlayerUsecase) Delete(ctx context.Context, id string) error {
	return uc.DBWriter.DeletePlayer(context.Background(), id)
}
