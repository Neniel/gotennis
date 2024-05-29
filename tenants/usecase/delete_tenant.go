package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"
)

type DeleteTenant interface {
	Do(ctx context.Context, customerID string) error
}

type deleteTenant struct {
	DBWriter database.DBWriter
}

func NewDeleteTenant(app app.IApp) DeleteTenant {
	systemMongoDBClient := app.GetSystemMongoDBClient()
	return &deleteTenant{
		DBWriter: database.NewDatabaseWriter(systemMongoDBClient.MongoDBClient, systemMongoDBClient.DatabaseName),
	}
}

func (uc *deleteTenant) Do(ctx context.Context, customerID string) error {

	err := uc.DBWriter.DeleteTenant(ctx, customerID)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not delete tenant: %w", err).Error())
		return err
	}
	return nil
}
