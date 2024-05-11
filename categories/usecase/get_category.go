package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
)

type GetCategoryUsecase interface {
	Get(ctx context.Context, id string) (*entity.Category, error)
}

type getCategoryUsecase struct {
	DBReader database.DBReader
}

func NewGetCategoryUsecase(app *app.App) GetCategoryUsecase {
	return &getCategoryUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *getCategoryUsecase) Get(ctx context.Context, id string) (*entity.Category, error) {
	return uc.DBReader.GetCategory(ctx, id)
}
