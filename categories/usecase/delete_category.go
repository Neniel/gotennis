package usecase

import (
	"context"
	"fmt"
	"log"

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
	err := uc.DBWriter.DeleteCategory(ctx, id)
	if err != nil {
		log.Println(fmt.Errorf("error at DeleteCategory: %w", err))
		return err
	}
	return nil
}
