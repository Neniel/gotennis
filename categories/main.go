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
			CreateCategoryUsecase: usecase.NewCreateCategory(dbWriter),
			DeleteCategory:        usecase.NewDeleteCategory(dbWriter),
			ListCategories:        usecase.NewListCategories(dbReader),
			GetCategory:           usecase.NewGetCategory(dbReader),
			UpdateCategory:        usecase.NewUpdateCategory(dbReader, dbWriter),
		},
	}

	ms.NewAPIServer().Run()
}
