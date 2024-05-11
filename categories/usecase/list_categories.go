package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
)

type ListCategoriesUsecase interface {
	List(ctx context.Context) ([]entity.Category, error)
}

type listCategoriesUsecase struct {
	DBReader database.DBReader
}

func NewListCategoriesUsecase(app *app.App) ListCategoriesUsecase {
	return &listCategoriesUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *listCategoriesUsecase) List(ctx context.Context) ([]entity.Category, error) {
	return uc.DBReader.GetCategories(ctx)
}
