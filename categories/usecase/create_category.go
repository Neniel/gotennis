package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/util"
)

type CreateCategoryUsecase interface {
	CreateCategory(ctx context.Context, request *CreateCategoryRequest) (*entity.Category, error)
}

type createCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewCreateCategoryUsecase(dbWriter database.DBWriter) CreateCategoryUsecase {
	return &createCategoryUsecase{
		DBWriter: dbWriter,
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (r *CreateCategoryRequest) Validate() error {
	if r.Name == "" {
		return util.ErrCategoryNameIsEmpty
	}

	return nil
}

func (uc *createCategoryUsecase) CreateCategory(ctx context.Context, request *CreateCategoryRequest) (*entity.Category, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not create category: %w", err).Error())
		return nil, err
	}

	category := entity.NewCategory(request.Name)

	newCategory, err := uc.DBWriter.AddCategory(ctx, category)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create category: %w", err).Error())
		return nil, err
	}
	return newCategory, nil
}
