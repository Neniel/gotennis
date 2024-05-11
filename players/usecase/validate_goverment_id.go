package usecase

import (
	"context"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/database"
)

type ValidateGovernmentIDUsecase interface {
	IsAvailable(ctx context.Context, id string) (bool, error)
}

type validateGovernmentIDUsecase struct {
	DBReader database.DBReader
}

func NewValidateGovernmentIDUsecase(app *app.App) ValidateGovernmentIDUsecase {
	return &validateGovernmentIDUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *validateGovernmentIDUsecase) IsAvailable(ctx context.Context, governmentID string) (bool, error) {
	return uc.DBReader.IsAvailable(ctx, "government_id", governmentID)
}
