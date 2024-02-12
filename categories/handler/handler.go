package handler

import (
	"categories/database"
	"context"
	"log"

	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"
)

type CategoriesHandler struct {
	CacheReader    database.DBReader
	CacheWriter    database.DBWriter
	DatabaseReader database.DBReader
	DatabaseWriter database.DBWriter
}

func NewCategoriesHandler(cacheReader database.DBReader, cacheWriter database.DBWriter, dbr database.DBReader, dbw database.DBWriter) *CategoriesHandler {
	return &CategoriesHandler{
		CacheReader:    cacheReader,
		CacheWriter:    cacheWriter,
		DatabaseReader: dbr,
		DatabaseWriter: dbw,
	}
}

func (h *CategoriesHandler) Add(ctx context.Context, categoryToAdd *entity.Category) (*entity.Category, error) {
	if len(categoryToAdd.Name) == 0 {
		return nil, util.ErrCategoryNameIsEmpty
	}

	return h.DatabaseWriter.AddCategory(ctx, categoryToAdd)
}

func (h *CategoriesHandler) GetAll(ctx context.Context) ([]entity.Category, error) {
	categories, err := h.CacheReader.GetCategories(ctx)
	if len(categories) == 0 || err != nil {
		if err != nil {
			log.Println(err.Error())
		}
		return h.DatabaseReader.GetCategories(ctx)
	}
	return categories, nil
}

func (h *CategoriesHandler) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	category, err := h.CacheReader.GetCategory(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return h.DatabaseReader.GetCategory(ctx, id)
	}
	return category, nil
}

func (h *CategoriesHandler) Delete(ctx context.Context, id string) error {
	return h.DatabaseWriter.DeleteCategory(ctx, id)
}
