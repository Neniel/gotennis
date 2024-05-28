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

type UpdateCustomer interface {
	Do(ctx context.Context, customerID string, request *UpdateCustomerRequest) (*entity.Customer, error)
}

type updateCustomer struct {
	DBWriter database.DBWriter
}

func NewUpdateCustomer(app app.IApp) UpdateCustomer {
	return &updateCustomer{
		DBWriter: database.NewDatabaseWriter(app.GetSystemMongoDBClient()),
	}
}

type UpdateCustomerRequest struct {
	Name string `json:"name"`
}

func (r *UpdateCustomerRequest) Validate() error {
	if r.Name == "" {
		return util.ErrCategoryNameIsEmpty
	}

	return nil
}

func (uc *updateCustomer) Do(ctx context.Context, customerID string, request *UpdateCustomerRequest) (*entity.Customer, error) {
	if err := request.Validate(); err != nil {
		log.Logger.Info(fmt.Errorf("could not update customer: %w", err).Error())
		return nil, err
	}

	customer := &entity.Customer{
		Name: request.Name,
	}

	newCustomer, err := uc.DBWriter.UpdateCustomer(ctx, customer)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not update customer: %w", err).Error())
		return nil, err
	}
	return newCustomer, nil
}
