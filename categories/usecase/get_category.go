package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"

	"github.com/Neniel/gotennis/lib/entity"
)

type GetCategory interface {
	Do(ctx context.Context, id string) (*entity.Category, error)
}

type getCategory struct {
	DBReader database.DBReader
}

func NewGetCategory(dbReader database.DBReader) GetCategory {
	return &getCategory{
		DBReader: dbReader,
	}
}

func (uc *getCategory) Do(ctx context.Context, id string) (*entity.Category, error) {
	category, err := uc.DBReader.GetCategory(ctx, id)
	if err != nil {
		log.Logger.Error(fmt.Errorf("could not get category: %w", err).Error())
		return nil, err
	}

	return category, nil
}
