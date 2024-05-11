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
			CreatePlayerUsecase: usecase.NewCreatePlayerUsecase(app),
			DeletePlayerUsecase: usecase.NewDeletePlayerUsecase(app),
			ListPlayersUsecase:  usecase.NewListPlayersUsecase(app),
			GetPlayerUsecase:    usecase.NewGetPlayerUsecase(app),
		},
	}

	ms.NewAPIServer().Run()
}
