package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"
)

type DeleteCategoryUsecase interface {
	Delete(ctx context.Context, id string) error
}

type deleteCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewDeleteCategoryUsecase(dbWriter database.DBWriter) DeleteCategoryUsecase {
	return &deleteCategoryUsecase{
		DBWriter: dbWriter,
	}
}

func (uc *deleteCategoryUsecase) Delete(ctx context.Context, id string) error {
	return uc.DBWriter.DeleteCategory(ctx, id)
}
