package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
)

type ListCustomers interface {
	Do(ctx context.Context) ([]entity.Customer, error)
}

type listCustomers struct {
	DBReader database.DBReader
}

func NewListCustomers(app app.IApp) ListCustomers {
	return &listCustomers{
		DBReader: database.NewDatabaseReader(app.GetSystemMongoDBClient()),
	}
}

func (uc *listCustomers) Do(ctx context.Context) ([]entity.Customer, error) {
	newCustomer, err := uc.DBReader.GetCustomers(ctx)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return newCustomer, nil
}
