package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
)

type ListTenants interface {
	Do(ctx context.Context) ([]entity.Tenant, error)
}

type listTenants struct {
	DBReader database.DBReader
}

func NewListTenants(app app.IApp) ListTenants {
	systemMongoDBClient := app.GetSystemMongoDBClient()
	return &listTenants{
		DBReader: database.NewDatabaseReader(systemMongoDBClient.MongoDBClient, systemMongoDBClient.DatabaseName),
	}
}

func (uc *listTenants) Do(ctx context.Context) ([]entity.Tenant, error) {
	newCustomer, err := uc.DBReader.GetTenants(ctx)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return newCustomer, nil
}
