package main

import (
	"categories/usecase"
	"context"

	"github.com/Neniel/gotennis/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &CategoryMicroservice{
		App: app,
		Usecases: &Usecases{
			CreateCategoryUsecase: usecase.NewCreateCategoryUsecase(app),
			SaveCategoryUsecase:   usecase.NewSaveCategoryUsecase(app),
			DeleteCategory:        usecase.NewDeleteCategoryUsecase(app),
			ListCategories:        usecase.NewListCategoriesUsecase(app),
			GetCategory:           usecase.NewGetCategoryUsecase(app),
		},
	}

	ms.NewAPIServer().Run()
}
