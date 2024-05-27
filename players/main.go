package main

import (
	"context"

	"github.com/Neniel/gotennis/players/usecase"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
)

func main() {
	app := app.NewApp(context.Background())
	dbReader := database.NewDatabaseReader(app.GetMongoDBClient())
	dbWriter := database.NewDatabaseWriter(app.GetMongoDBClient())

	ms := &PlayerMicroservice{
		App: app,
		Usecases: &Usecases{
			CreatePlayer: usecase.NewCreatePlayer(dbWriter, dbReader),
			DeletePlayer: usecase.NewDeletePlayer(dbWriter),
			ListPlayers:  usecase.NewListPlayers(dbReader),
			GetPlayer:    usecase.NewGetPlayer(dbReader),
			UpdatePlayer: usecase.NewUpdatePlayer(dbWriter, dbReader),
		},
	}

	ms.NewAPIServer().Run()
}
