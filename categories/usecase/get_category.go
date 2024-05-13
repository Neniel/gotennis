package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
)

type GetCategoryUsecase interface {
	Get(ctx context.Context, id string) (*entity.Category, error)
}

type getCategoryUsecase struct {
	DBReader database.DBReader
}

func NewGetCategoryUsecase(dbReader database.DBReader) GetCategoryUsecase {
	return &getCategoryUsecase{
		DBReader: dbReader,
	}
}

func (uc *getCategoryUsecase) Get(ctx context.Context, id string) (*entity.Category, error) {
	return uc.DBReader.GetCategory(ctx, id)
}
