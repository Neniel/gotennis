package main

import (
	"context"

	"github.com/Neniel/gotennis/tournaments/usecase"

	"github.com/Neniel/gotennis/lib/app"
	"github.com/Neniel/gotennis/lib/database"
)

func main() {
	app := app.NewApp(context.Background())
	dbReader := database.NewDatabaseReader(app.GetMongoDBClient())
	dbWriter := database.NewDatabaseWriter(app.GetMongoDBClient())

	ms := &TournamentMicroservice{
		App: app,
		Usecases: &Usecases{
			CreateTournament: usecase.NewCreateTournament(dbWriter),
			DeleteTournament: usecase.NewDeleteTournament(dbWriter),
			ListTournaments:  usecase.NewListTournaments(dbReader),
			GetTournament:    usecase.NewGetTournament(dbReader),
			UpdateTournament: usecase.NewUpdateTournament(dbWriter, dbReader),
		},
	}

	ms.NewAPIServer().Run()
}
