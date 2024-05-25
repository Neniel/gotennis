package usecase

import (
	"context"
	"errors"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"
)

type UpdateCategory interface {
	Do(ctx context.Context, id string, request *UpdateCategoryRequest) (*entity.Category, error)
}

type updateCategory struct {
	DBReader database.DBReader
	DBWriter database.DBWriter
}

func NewUpdateCategory(dbReader database.DBReader, dbWriter database.DBWriter) UpdateCategory {
	return &updateCategory{
		DBReader: dbReader,
		DBWriter: dbWriter,
	}
}

type UpdateCategoryRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r *UpdateCategoryRequest) Validate(id string) error {

	if r.ID == "" {
		return errors.New("category ID is required for update")
	}

	if r.ID != id {
		return errors.New("provided category ID does not match the ID of the category to be updated")
	}

	if r.Name == "" {
		return util.ErrCategoryNameIsEmpty
	}

	return nil
}

func (uc *updateCategory) Do(ctx context.Context, id string, request *UpdateCategoryRequest) (*entity.Category, error) {
	if err := request.Validate(id); err != nil {
		return nil, err
	}

	category, err := uc.DBReader.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = request.Name

	updatedCategory, err := uc.DBWriter.UpdateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return updatedCategory, err
}
