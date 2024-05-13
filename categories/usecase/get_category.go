package usecase

import (
	"context"
	"fmt"
	"log"

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
	category, err := uc.DBReader.GetCategory(ctx, id)
	if err != nil {
		log.Println(fmt.Errorf("error at GetCategory: %w", err))
		return nil, err
	}
	return category, nil
}
