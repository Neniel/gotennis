package usecase

import (
	"github.com/Neniel/gotennis/entity"
)

type ValidateCategoryUsecase interface {
	ValidateCategory(category *entity.Category) error
}

type validateCategoryUsecase struct{}

func NewValidateCategoryUsecase() ValidateCategoryUsecase {
	return &validateCategoryUsecase{}
}

func (uc *validateCategoryUsecase) ValidateCategory(category *entity.Category) error {
	//return category.Validate()
	return nil
}
