package main

import (
	"context"

	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &CategoryMicroservice{
		App: app,
		//Usecases: &Usecases{
		//	CreateCategoryUsecase: usecase.NewCreateCategory(dbWriter),
		//	DeleteCategory:        usecase.NewDeleteCategory(dbWriter),
		//	ListCategories:        usecase.NewListCategories(dbReader),
		//	GetCategory:           usecase.NewGetCategory(dbReader),
		//	UpdateCategory:        usecase.NewUpdateCategory(dbReader, dbWriter),
		//},
	}

	ms.NewAPIServer().Run()
}
