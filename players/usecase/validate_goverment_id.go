package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type ValidateGovernmentIDUsecase interface {
	IsAvailable(ctx context.Context, governmentID string) (bool, error)
}

type validateGovernmentIDUsecase struct {
	DBReader database.DBReader
}

func NewValidateGovernmentIDUsecase(dbReader database.DBReader) ValidateGovernmentIDUsecase {
	return &validateGovernmentIDUsecase{
		DBReader: dbReader,
	}
}

func (uc *validateGovernmentIDUsecase) IsAvailable(ctx context.Context, governmentID string) (bool, error) {
	if governmentID == "" {
		return false, nil
	}

	return uc.DBReader.IsAvailable(ctx, "government_id", governmentID)
}
