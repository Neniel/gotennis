package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
)

type SaveCategoryUsecase interface {
	Save(ctx context.Context, category *entity.Category) (*entity.Category, error)
}

type saveCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewSaveCategoryUsecase(app *app.App) SaveCategoryUsecase {
	return &saveCategoryUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
	}
}

func (uc *saveCategoryUsecase) Save(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return uc.DBWriter.AddCategory(ctx, category)
}
