package main

import (
	"context"
	"players/usecase"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/database"
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
