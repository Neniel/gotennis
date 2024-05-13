package usecase

import (
	"context"

	"github.com/Neniel/gotennis/database"
)

type ValidateEmailUsecase interface {
	IsAvailable(ctx context.Context, email string) (bool, error)
}

type validateEmailUsecaseUsecase struct {
	DBReader database.DBReader
}

func NewValidateEmailUsecase(dbReader database.DBReader) ValidateEmailUsecase {
	return &validateEmailUsecaseUsecase{
		DBReader: dbReader,
	}
}

func (uc *validateEmailUsecaseUsecase) IsAvailable(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, nil
	}

	return uc.DBReader.IsAvailable(ctx, "email", email)
}
