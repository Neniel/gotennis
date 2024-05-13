package usecase

import (
	"context"

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
	return uc.DBReader.GetCategories(ctx)
}
