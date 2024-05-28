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

type CreateCustomer interface {
	Do(ctx context.Context, request *CreateCustomerRequest) (*entity.Category, error)
}

type createCategory struct {
	App      app.IApp
	DBWriter database.DBWriter
}

func NewCreateCustomer(app app.IApp) (CreateCustomer, error) {
	return &createCategory{
		DBWriter: database.NewDatabaseWriter(mongoDBClient),
	}, nil
}

type CreateCustomerRequest struct {
	Name string `json:"name"`
}

func (r *CreateCustomerRequest) Validate() error {
	if r.Name == "" {
		return util.ErrCategoryNameIsEmpty
	}

	return nil
}

func (uc *createCategory) Do(ctx context.Context, request *CreateCustomerRequest) (*entity.Category, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}

	customer := entity.Customer{}

	newCategory, err := uc.DBWriter.AddCategory(ctx, customer)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return newCategory, nil
}
