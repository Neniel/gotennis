package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Neniel/gotennis/database"

	"github.com/Neniel/gotennis/entity"
)

type ListCategoriesUsecase interface {
	List(ctx context.Context) ([]entity.Category, error)
}

type listCategoriesUsecase struct {
	DBReader database.DBReader
}

func NewListCategoriesUsecase(dbReader database.DBReader) ListCategoriesUsecase {
	return &listCategoriesUsecase{
		DBReader: dbReader,
	}
}

func (uc *listCategoriesUsecase) List(ctx context.Context) ([]entity.Category, error) {
	categories, err := uc.DBReader.GetCategories(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error at GetCategories: %w", err))
		return nil, err
	}
	return categories, nil
}
