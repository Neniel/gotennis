package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"

	"github.com/Neniel/gotennis/lib/entity"
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
	category, err := uc.DBReader.GetCategory(ctx, id)
	if err != nil {
		log.Logger.Error(fmt.Errorf("could not get category: %w", err).Error())
		return nil, err
	}

	return category, nil
}
