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
			CreatePlayerUsecase: usecase.NewCreatePlayerUsecase(dbWriter, dbReader),
			DeletePlayerUsecase: usecase.NewDeletePlayerUsecase(dbWriter),
			ListPlayersUsecase:  usecase.NewListPlayersUsecase(dbReader),
			GetPlayerUsecase:    usecase.NewGetPlayerUsecase(dbReader),
			UpdatePlayerUsecase: usecase.NewUpdatePlayerUsecase(dbWriter, dbReader),
		},
	}

	ms.NewAPIServer().Run()
}
