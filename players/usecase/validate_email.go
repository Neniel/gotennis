package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type ValidateEmail interface {
	IsAvailable(ctx context.Context, email string) (bool, error)
}

type validateEmailUsecaseUsecase struct {
	DBReader database.DBReader
}

func NewValidateEmailUsecase(dbReader database.DBReader) ValidateEmail {
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
