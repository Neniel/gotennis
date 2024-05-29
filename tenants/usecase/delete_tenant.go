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
	return &deleteTenant{
		DBWriter: database.NewDatabaseWriter(app.GetSystemMongoDBClient(), "system"),
	}
}

func (uc *deleteTenant) Do(ctx context.Context, customerID string) error {

	err := uc.DBWriter.DeleteTournament(ctx, customerID)
	if err != nil {
		log.Logger.Info(fmt.Errorf("could not update customer: %w", err).Error())
		return err
	}
	return nil
}
