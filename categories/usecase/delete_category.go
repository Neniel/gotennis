package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Neniel/gotennis/lib/database"
)

type DeleteCategory interface {
	Do(ctx context.Context, id string) error
}

type deleteCategory struct {
	DBWriter database.DBWriter
}

func NewDeleteCategory(dbWriter database.DBWriter) DeleteCategory {
	return &deleteCategory{
		DBWriter: dbWriter,
	}
}

func (uc *deleteCategory) Do(ctx context.Context, id string) error {
	err := uc.DBWriter.DeleteCategory(ctx, id)
	if err != nil {
		log.Println(fmt.Errorf("error at DeleteCategory: %w", err))
		return err
	}
	return nil
}
