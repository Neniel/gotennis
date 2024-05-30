package main

import (
	"context"

	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &AuthMicroservice{
		App: app,
		/*
			Usecases: &Usecases{
				CreateTournament: usecase.NewCreateTournament(dbWriter),
				DeleteTournament: usecase.NewDeleteTournament(dbWriter),
				ListTournaments:  usecase.NewListTournaments(dbReader),
				GetTournament:    usecase.NewGetTournament(dbReader),
				UpdateTournament: usecase.NewUpdateTournament(dbWriter, dbReader),
			},
		*/
	}

	ms.NewAPIServer().Run()
}
