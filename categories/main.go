package main

import (
	"context"

	"github.com/Neniel/gotennis/categories/usecase"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
)

func main() {
	app := app.NewApp(context.Background())
	dbReader := database.NewDatabaseReader(app.GetMongoDBClient())
	dbWriter := database.NewDatabaseWriter(app.GetMongoDBClient())

	ms := &CategoryMicroservice{
		App: app,
		Usecases: &Usecases{
			CreateCategoryUsecase: usecase.NewCreateCategoryUsecase(dbWriter),
			DeleteCategory:        usecase.NewDeleteCategoryUsecase(dbWriter),
			ListCategories:        usecase.NewListCategoriesUsecase(dbReader),
			GetCategory:           usecase.NewGetCategoryUsecase(dbReader),
			UpdateCategory:        usecase.NewUpdateCategoryUsecase(dbReader, dbWriter),
		},
	}

	ms.NewAPIServer().Run()
}
