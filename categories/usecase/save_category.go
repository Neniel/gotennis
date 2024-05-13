package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
)

type SaveCategoryUsecase interface {
	Save(ctx context.Context, category *entity.Category) (*entity.Category, error)
}

type saveCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewSaveCategoryUsecase(dbWriter database.DBWriter) SaveCategoryUsecase {
	return &saveCategoryUsecase{
		DBWriter: dbWriter,
	}
}

func (uc *saveCategoryUsecase) Save(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return uc.DBWriter.AddCategory(ctx, category)
}
