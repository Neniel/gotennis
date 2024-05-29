package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
)

type GetCustomer interface {
	Do(ctx context.Context, id string) (*entity.Customer, error)
}

type getCustomer struct {
	DBReader database.DBReader
}

func NewGetCustomer(app app.IApp) GetCustomer {
	return &getCustomer{
		DBReader: database.NewDatabaseReader(app.GetSystemMongoDBClient()),
	}
}

func (uc *getCustomer) Do(ctx context.Context, id string) (*entity.Customer, error) {
	customer, err := uc.DBReader.GetCustomer(ctx, id)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return customer, nil
}
