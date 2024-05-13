package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type CreateCategoryUsecase interface {
	CreateCategory(ctx context.Context, request *CreateCategoryRequest) (*entity.Category, error)
}

type createCategoryUsecase struct {
	SaveCategoryUsecase SaveCategoryUsecase
}

func NewCreateCategoryUsecase(dbWriter database.DBWriter) CreateCategoryUsecase {
	return &createCategoryUsecase{
		SaveCategoryUsecase: NewSaveCategoryUsecase(dbWriter),
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
		return nil, err
	}

	category := entity.NewCategory(request.Name)

	return uc.SaveCategoryUsecase.Save(ctx, category)
}
