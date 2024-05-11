package usecase

import (
	"context"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type CreateCategoryUsecase interface {
	CreateCategory(request *CreateCategoryRequest) (*entity.Category, error)
}

type createCategoryUsecase struct {
	DBWriter database.DBWriter
}

func NewCreateCategoryUsecase(app *app.App) CreateCategoryUsecase {
	return &createCategoryUsecase{
		DBWriter: database.NewDatabaseWriter(app.DBClients.MongoDB),
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

func (uc *createCategoryUsecase) CreateCategory(request *CreateCategoryRequest) (*entity.Category, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	category := entity.NewCategory(request.Name)

	return uc.DBWriter.AddCategory(context.Background(), category)
}
