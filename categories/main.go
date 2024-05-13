package main

import (
	"categories/usecase"
	"context"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/database"
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
