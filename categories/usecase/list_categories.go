package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Neniel/gotennis/lib/database"

	"github.com/Neniel/gotennis/lib/entity"
)

type ListCategories interface {
	Do(ctx context.Context) ([]entity.Category, error)
}

type listCategories struct {
	DBReader database.DBReader
}

func NewListCategories(dbReader database.DBReader) ListCategories {
	return &listCategories{
		DBReader: dbReader,
	}
}

func (uc *listCategories) Do(ctx context.Context) ([]entity.Category, error) {
	categories, err := uc.DBReader.GetCategories(ctx)
	if err != nil {
		log.Println(fmt.Errorf("error at GetCategories: %w", err))
		return nil, err
	}
	return categories, nil
}
