package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"github.com/Neniel/gotennis/lib/util"
)

type CreateCategory interface {
	CreateCategory(ctx context.Context, request *CreateCategoryRequest) (*entity.Category, error)
}

type createCategory struct {
	App      app.IApp
	DBWriter database.DBWriter
}

func NewCreateCategory(app app.IApp, customerID string) (CreateCategory, error) {
	mongoDBClient, ok := app.GetMongoDBClients()[customerID]
	if !ok {
		log.Logger.Info(fmt.Errorf("could not initialize usecase CreateCategory: no mongoDBClient for customerID '%s'", customerID).Error())
		return nil, fmt.Errorf("could not initialize usecase CreateCategory: no mongoDBClient for customerID '%s'", customerID)
	}

	return &createCategory{
		DBWriter: database.NewDatabaseWriter(mongoDBClient),
	}, nil
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

func (uc *createCategory) CreateCategory(ctx context.Context, request *CreateCategoryRequest) (*entity.Category, error) {
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
