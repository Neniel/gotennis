package main

import (
	"context"
	"players/usecase"

	"github.com/Neniel/gotennis/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &PlayerMicroservice{
		App: app,
		Usecases: &Usecases{
			CreatePlayerUsecase: usecase.NewCreateCategoryUsecase(app),
			SavePlayerUsecase:   nil,
			DeletePlayerUsecase: nil,
			ListPlayersUsecase:  nil,
			GetPlayerUsecase:    nil,
		},
	}

	ms.NewAPIServer().Run()
}
