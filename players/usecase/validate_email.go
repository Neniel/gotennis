package usecase

import (
	"context"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/database"
)

type ValidateEmailUsecase interface {
	IsAvailable(ctx context.Context, email string) (bool, error)
}

type validateEmailUsecaseUsecase struct {
	DBReader database.DBReader
}

func NewValidateEmailUsecase(app *app.App) ValidateEmailUsecase {
	return &validateEmailUsecaseUsecase{
		DBReader: database.NewDatabaseReader(app.DBClients.MongoDB),
	}
}

func (uc *validateEmailUsecaseUsecase) IsAvailable(ctx context.Context, email string) (bool, error) {
	return uc.DBReader.IsAvailable(ctx, "email", email)
}
