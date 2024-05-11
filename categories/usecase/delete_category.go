package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
)

type DeleteCategoryUsecase interface {
	Delete(ctx context.Context, id string) error
}

type deleteCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewDeleteCategoryUsecase(app *app.App) DeleteCategoryUsecase {
	return &deleteCategoryUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
	}
}

func (uc *deleteCategoryUsecase) Delete(ctx context.Context, id string) error {
	return uc.DBWriter.DeleteCategory(ctx, id)
}
