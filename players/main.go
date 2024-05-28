package main

import (
	"context"

	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &PlayerMicroservice{
		App: app,
		/*
			Usecases: &Usecases{
				CreatePlayer: usecase.NewCreatePlayer(dbWriter, dbReader),
				DeletePlayer: usecase.NewDeletePlayer(dbWriter),
				ListPlayers:  usecase.NewListPlayers(dbReader),
				GetPlayer:    usecase.NewGetPlayer(dbReader),
				UpdatePlayer: usecase.NewUpdatePlayer(dbWriter, dbReader),
			},
		*/
	}

	ms.NewAPIServer().Run()
}
