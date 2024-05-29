package usecase

import (
	"context"
	"fmt"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/log"
)

type DeleteCustomer interface {
	Do(ctx context.Context, customerID string) error
}

type deleteCustomer struct {
	DBWriter database.DBWriter
}

func NewDeleteCustomer(app app.IApp) DeleteCustomer {
	return &deleteCustomer{
		DBWriter: database.NewDatabaseWriter(app.GetSystemMongoDBClient()),
	}
}

func (uc *deleteCustomer) Do(ctx context.Context, customerID string) error {

	err := uc.DBWriter.DeleteCustomer(ctx, customerID)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not update customer: %w", err).Error())
		return err
	}
	return nil
}
