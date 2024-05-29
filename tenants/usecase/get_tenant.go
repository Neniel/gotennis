package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
)

type GetTenant interface {
	Do(ctx context.Context, id string) (*entity.Tenant, error)
}

type getTenant struct {
	DBReader database.DBReader
}

func NewGetTenant(app app.IApp) GetTenant {
	systemMongoDBClient := app.GetSystemMongoDBClient()
	return &getTenant{
		DBReader: database.NewDatabaseReader(systemMongoDBClient.MongoDBClient, systemMongoDBClient.DatabaseName),
	}
}

func (uc *getTenant) Do(ctx context.Context, id string) (*entity.Tenant, error) {
	customer, err := uc.DBReader.GetTenant(ctx, id)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not create customer: %w", err).Error())
		return nil, err
	}
	return customer, nil
}
