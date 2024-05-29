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
	Do(ctx context.Context, request *CreateCustomerRequest) (*entity.Customer, error)
}

type createCustomer struct {
	DBWriter database.DBWriter
}

func NewCreateCustomer(app app.IApp) CreateCustomer {
	return &createCustomer{
		DBWriter: database.NewDatabaseWriter(app.GetSystemMongoDBClient()),
	}
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

func (uc *createCustomer) Do(ctx context.Context, request *CreateCustomerRequest) (*entity.Customer, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}

	customer := &entity.Customer{
		Name: request.Name,
	}

	newCustomer, err := uc.DBWriter.AddCustomer(ctx, customer)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return newCustomer, nil
}
