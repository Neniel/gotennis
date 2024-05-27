package usecase

import (
	"context"

	"github.com/Neniel/gotennis/lib/database"
)

type ValidateAlias interface {
	IsAvailable(ctx context.Context, alias *string) (bool, error)
}

type validateAliasUsecase struct {
	DBReader database.DBReader
}

func NewValidateAliasUsecase(dbReader database.DBReader) ValidateAlias {
	return &validateAliasUsecase{
		DBReader: dbReader,
	}
}

func (uc *validateAliasUsecase) IsAvailable(ctx context.Context, alias *string) (bool, error) {
	if alias == nil {
		return true, nil
	}

	return uc.DBReader.IsAvailable(ctx, "alias", *alias)
}
